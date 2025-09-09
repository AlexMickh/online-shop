package token_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) SaveToken(ctx context.Context, userId, token, tokenType string) error {
	const op = "repository.postgres.token.SaveToken"

	query := `INSERT INTO tokens
			  (user_id, token, type)
			  VALUES ($1, $2, $3)`
	_, err := p.db.Exec(ctx, query, userId, token, tokenType)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *Postgres) UserIdByToken(ctx context.Context, token, tokenType string) (string, error) {
	const op = "repository.postgres.token.UserIdByToken"

	query := "SELECT user_id FROM tokens WHERE token = $1 AND type = $2"
	var userId string
	err := p.db.QueryRow(ctx, query, token, tokenType).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, errs.ErrTokenNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}
