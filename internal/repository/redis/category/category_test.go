package category_cash

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestCash_SaveCategory(t *testing.T) {
	type fields struct {
		rdb    *redis.Client
		expire time.Duration
	}
	type args struct {
		ctx      context.Context
		categoty models.Category
	}

	rdb := initCash(t)
	defer rdb.Close()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "good case",
			fields: fields{
				rdb:    rdb,
				expire: time.Hour,
			},
			args: args{
				ctx: context.Background(),
				categoty: models.Category{
					ID:   uuid.NewString(),
					Name: "phones",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cash{
				rdb:    tt.fields.rdb,
				expire: tt.fields.expire,
			}
			if err := c.SaveCategory(tt.args.ctx, tt.args.categoty); err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Cash.SaveCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	for _, tt := range tests {
		_ = rdb.Del(tt.args.ctx, genKey(tt.args.categoty.ID)).Err()
	}
}

func TestCash_AllCategories(t *testing.T) {
	type fields struct {
		rdb    *redis.Client
		expire time.Duration
	}
	type args struct {
		ctx context.Context
	}

	rdb := initCash(t)
	defer rdb.Close()

	categories := make([]models.Category, 10)
	for i := range 10 {
		categories[i] = models.Category{
			ID:   uuid.NewString(),
			Name: gofakeit.BeerStyle(),
		}
		err := rdb.HSet(context.Background(), genKey(categories[i].ID), categories[i]).Err()
		if err != nil {
			t.Fatalf("failed to save: %v, index: %d", err, i)
		}
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
				rdb:    rdb,
				expire: time.Hour,
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
			c := &Cash{
				rdb:    tt.fields.rdb,
				expire: tt.fields.expire,
			}
			got, err := c.AllCategories(tt.args.ctx)
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Cash.AllCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == len(tt.want) {
				t.Errorf("Cash.AllCategories() = %v, want %v", got, tt.want)
			}
		})
	}
	for i, tt := range tests {
		_ = rdb.Del(tt.args.ctx, genKey(tt.want[i].ID)).Err()
	}
}

func initCash(t *testing.T) *redis.Client {
	t.Helper()

	// db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	// if err != nil {
	// 	t.Fatalf("failed to init cash: %v", err)
	// }

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	// 	Password: os.Getenv("REDIS_PASSWORD"),
	// 	DB:       db,
	// })

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6359",
		Password: "root",
		DB:       0,
	})

	return rdb
}
