package postgres

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	commonErrors "github.com/Forvi/MayWalk/internal/common/errors"
	"github.com/Forvi/MayWalk/internal/user/adapters/repository/models"
	"github.com/Forvi/MayWalk/internal/user/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEmptyUpdates = commonErrors.NewDomainError("internal.empty_updates", "no updates fields", http.StatusInternalServerError)

type querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return repo.db.Begin(ctx)
}

func (repo *UserRepo) SaveWithTx(ctx context.Context, user domain.User, tx pgx.Tx) (*domain.User, error) {
	exists, err := repo.findByIDWithQuerier(ctx, user.ID, tx)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
	}
	if exists != nil {
		return repo.update(ctx, &user, exists, tx)
	}
	return repo.create(ctx, &user, tx)
}

func (repo *UserRepo) Save(ctx context.Context, user domain.User) (*domain.User, error) {
	exists, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
	}
	if exists != nil {
		return repo.update(ctx, &user, exists, repo.db)
	}
	return repo.create(ctx, &user, repo.db)
}

func (repo *UserRepo) FindByName(ctx context.Context, name domain.UserName) (*domain.User, error) {
	query := `
	SELECT id, username, email, password_hash, is_active, created_at, updated_at 
	FROM users 
	WHERE username=$1 AND is_active=true;
	`

	rows, err := repo.db.Query(ctx, query, &name)
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}
	defer rows.Close()

	userRow, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.DBUser])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("collect user: %w", err)
	}

	return userRow.ToDomain()
}

func (repo *UserRepo) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	return repo.findByIDWithQuerier(ctx, id, repo.db)
}

func (repo *UserRepo) FindByEmail(ctx context.Context, email domain.UserEmail) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, is_active, created_at, updated_at 
              FROM users 
              WHERE email=$1 AND is_active=true`

	rows, err := repo.db.Query(ctx, query, &email)
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}
	defer rows.Close()

	userRow, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.DBUser])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("collect user: %w", err)
	}

	return userRow.ToDomain()
}

func (repo *UserRepo) IsUsernameAvailable(ctx context.Context, name domain.UserName) (bool, error) {
	query := `SELECT NOT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var available bool
	err := repo.db.QueryRow(ctx, query, name).Scan(&available)
	if err != nil {
		return false, fmt.Errorf("check username availability: %w", err)
	}

	return available, nil
}

func (repo *UserRepo) findByIDWithQuerier(ctx context.Context, id domain.UserID, q querier) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, is_active, created_at, updated_at FROM users WHERE id=$1`

	rows, err := q.Query(ctx, query, &id)
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}
	defer rows.Close()

	userRow, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.DBUser])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("collect user: %w", err)
	}

	return userRow.ToDomain()
}

func (repo *UserRepo) create(ctx context.Context, user *domain.User, q querier) (*domain.User, error) {
	query := `INSERT INTO users (id, username, email, password_hash, is_active, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
     		  ON CONFLICT (id) DO NOTHING;`

	result, err := q.Exec(ctx, query,
		user.ID, user.Name, user.Email, user.Password,
		user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return nil, domain.ErrAlreadyExists
	}

	return user, nil
}

func (repo *UserRepo) update(ctx context.Context, user, exists *domain.User, q querier) (*domain.User, error) {
	updates := repo.diff(user, exists)
	if len(updates) == 1 {
		return nil, ErrEmptyUpdates
	}
	setClauses := make([]string, 0, len(updates))
	args := make([]any, 0, len(updates)+1)
	i := 1

	for col, val := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}
	args = append(args, exists.ID)

	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), i,
	)

	_, err := q.Exec(ctx, query, args...)
	return user, err
}

func (repo *UserRepo) diff(user, exists *domain.User) map[string]any {
	updates := make(map[string]any)
	if exists.Name != user.Name {
		updates["username"] = user.Name
	}
	if exists.Email != user.Email {
		updates["email"] = user.Email
	}
	if exists.IsActive != user.IsActive {
		updates["is_active"] = user.IsActive
	}
	updates["updated_at"] = time.Now().UTC()
	return updates
}
