package product_repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgres_SaveProduct(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx         context.Context
		productId   string
		name        string
		description string
		price       float32
		imageUrl    string
		categoryIds []string
	}

	pool := initStorage()
	defer pool.Close()

	catigories := []models.Category{
		{
			ID:   uuid.NewString(),
			Name: gofakeit.CarType(),
		},
		{
			ID:   uuid.NewString(),
			Name: gofakeit.CarType(),
		},
	}
	categoryIds := make([]string, 0, len(catigories))

	for _, category := range catigories {
		_, _ = pool.Exec(context.Background(), "INSERT INTO categories (id, name) VALUES ($1, $2)", category.ID, category.Name)
		categoryIds = append(categoryIds, category.ID)
	}

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
				ctx:         context.Background(),
				productId:   uuid.NewString(),
				name:        "iphone",
				description: "gvdsvs",
				price:       567.8,
				imageUrl:    "bfdlknbvldvn",
				categoryIds: categoryIds,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				db: tt.fields.db,
			}
			if err := p.SaveProduct(
				tt.args.ctx,
				tt.args.productId,
				tt.args.name,
				tt.args.description,
				tt.args.price,
				tt.args.imageUrl,
				tt.args.categoryIds,
			); !errors.Is(err, tt.wantErr) {
				t.Errorf("Postgres.SaveProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	for _, tt := range tests {
		_, _ = pool.Exec(tt.args.ctx, "DELETE FROM products WHERE id = $1", tt.args.productId)
	}
	for _, categoryId := range categoryIds {
		_, _ = pool.Exec(context.Background(), "DELETE FROM categories WHERE id = $1", categoryId)
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
