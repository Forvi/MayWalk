package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Forvi/MayWalk/internal/user/adapters/repository/models"
	"github.com/Forvi/MayWalk/internal/user/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OAuthRepo struct {
	db *pgxpool.Pool
}

func NewOAuthRepo(db *pgxpool.Pool) *OAuthRepo {
	return &OAuthRepo{db: db}
}

func (repo *OAuthRepo) CreateWithTx(ctx context.Context, tx pgx.Tx, account models.OAuthAccountCreate) error {
	query := `INSERT INTO oauth_accounts (user_id, provider, provider_user_id) 
			  VALUES ($1, $2, $3);`

	_, err := tx.Exec(ctx, query, account.UserID, account.Provider, account.ProviderUserID)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.ErrAlreadyExists
			}
		}
		return fmt.Errorf("execute insert oauth account: %w", err)
	}
	return nil
}

func (repo *OAuthRepo) FindByProvider(ctx context.Context, provider, providerUserID string) (*models.OAuthAccount, error) {
	query := `SELECT * FROM oauth_accounts WHERE provider=$1 AND provider_user_id=$2;`

	rows, err := repo.db.Query(ctx, query, provider, providerUserID)
	if err != nil {
		return nil, fmt.Errorf("query oauth account: %w", err)
	}
	defer rows.Close()

	account, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.OAuthAccount])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("collect oauth account: %w", err)
	}

	return &account, nil
}

func (repo *OAuthRepo) FindByUserID(ctx context.Context, userID domain.UserID) ([]*models.OAuthAccount, error) {
	query := `SELECT * FROM oauth_accounts WHERE user_id=$1;`

	rows, err := repo.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query oauth account: %w", err)
	}
	defer rows.Close()

	accounts, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.OAuthAccount])
	if err != nil {
		return nil, fmt.Errorf("collect oauth accounts: %w", err)
	}

	return accounts, nil
}
