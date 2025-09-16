package category_repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestPostgres_SaveCategory(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx  context.Context
		id   string
		name string
	}

	pool := initStorage()
	defer pool.Close()

	id := uuid.NewString()
	name := "blkrgberkle"

	pool.Exec(context.Background(), "INSERT INTO categories (id, name) VALUES ($1, $2)", id, name)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "good case",
			fields: fields{
				db: pool,
			},
			args: args{
				ctx:  context.Background(),
				id:   uuid.NewString(),
				name: "klvdfsvlk",
			},
			wantErr: nil,
		},
		{
			name: "already exists case",
			fields: fields{
				db: pool,
			},
			args: args{
				ctx:  context.Background(),
				id:   id,
				name: name,
			},
			wantErr: errs.ErrCategoryAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				db: tt.fields.db,
			}
			if err := p.SaveCategory(tt.args.ctx, tt.args.id, tt.args.name); err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Postgres.SaveCategory() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
	for _, tt := range tests {
		_, _ = pool.Exec(tt.args.ctx, "DELETE FROM categories WHERE id = $1", tt.args.id)
	}
}

func TestPostgres_AllCategories(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
	}

	pool := initStorage()
	defer pool.Close()

	categories := make([]models.Category, 10)
	for i := range 10 {
		categories[i] = models.Category{
			ID:   uuid.NewString(),
			Name: gofakeit.BeerStyle(),
		}
		_, err := pool.Exec(
			context.Background(),
			"INSERT INTO categories (id, name) VALUES ($1, $2)",
			categories[i].ID,
			categories[i].Name,
		)
		require.NoError(t, err, "failed to save")
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Category
		wantErr error
	}{
		{
			name: "good case",
			fields: fields{
				db: pool,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    categories,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				db: tt.fields.db,
			}
			got, err := p.AllCategories(tt.args.ctx)
			require.ErrorIs(t, err, tt.wantErr, fmt.Sprintf("Postgres.AllCategories() error = %v, wantErr %v", err, tt.wantErr))

			fmt.Println(len(got))
			require.Equal(t, len(tt.want), len(got), "Postgres.AllCategories() = ", got)
		})
	}
	for _, tt := range tests {
		for _, w := range tt.want {
			_, _ = pool.Exec(tt.args.ctx, "DELETE FROM categories WHERE id = $1", w.ID)
		}
	}
}

func initStorage() *pgxpool.Pool {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=%s&pool_min_conns=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_MIN_POOLS"),
		os.Getenv("DB_MAX_POOLS"),
	)

	pool, _ := pgxpool.New(context.Background(), connString)

	return pool
}
