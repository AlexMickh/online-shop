package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	"github.com/jackc/pgx/v5/pgconn"
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

func (p *Postgres) SaveCategory(ctx context.Context, id string, name string) error {
	const op = "repository.postgres.category.SaveCategory"

	query := "INSERT INTO categories (id, name) VALUES ($1, $2)"
	_, err := p.db.Exec(ctx, query, id, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, errs.ErrCategoryAlreadyExists)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
