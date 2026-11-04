package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Forvi/MayWalk/internal/user/adapters/repository/models"
	"github.com/Forvi/MayWalk/internal/user/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) Save(ctx context.Context, user domain.User) error {
	exists, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return err
		}
	}
	if exists != nil {
		return repo.update(ctx, user, *exists)
	}
	return repo.create(ctx, user)
}

func (repo *UserRepo) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, is_active, created_at, updated_at FROM users WHERE id=$1`

	rows, err := repo.db.Query(ctx, query, &id)
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

func (repo *UserRepo) create(ctx context.Context, user domain.User) error {
	query := `
        INSERT INTO users (id, username, email, password_hash, is_active, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7);
    `
	_, err := repo.db.Exec(ctx, query,
		user.ID, user.Name, user.Email, user.Password,
		user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.ErrAlreadyExists
			}
		}
		return fmt.Errorf("execute insert user query: %w", err)
	}
	return nil
}

func (repo *UserRepo) update(ctx context.Context, user, exists domain.User) error {
	updates := repo.diff(user, exists)
	if len(updates) == 0 {
		return nil
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

	for _, a := range args {
		fmt.Printf("%v, ", a)
	}

	_, err := repo.db.Exec(ctx, query, args...)
	return err
}

func (repo *UserRepo) diff(user, exists domain.User) map[string]any {
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
