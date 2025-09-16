package category_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
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

func (p *Postgres) AllCategories(ctx context.Context) ([]models.Category, error) {
	const op = "repository.postgres.category.AllCategories"

	query := "SELECT id, name FROM categories"
	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		categories = append(categories, category)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return categories, nil
}
