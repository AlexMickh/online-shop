package session_cash

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestCash_SaveSession(t *testing.T) {
	type fields struct {
		rdb    *redis.Client
		expire time.Duration
	}
	type args struct {
		ctx  context.Context
		id   string
		user models.User
	}

	rdb := initCash()
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
				expire: time.Minute,
			},
			args: args{
				ctx: context.Background(),
				id:  uuid.NewString(),
				user: models.User{
					ID:              uuid.NewString(),
					Login:           "test",
					Email:           "test",
					Password:        "test",
					Role:            "test",
					IsEmailVerified: false,
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
			if err := c.SaveSession(tt.args.ctx, tt.args.id, tt.args.user); err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Cash.SaveSession() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func initCash() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6359",
		Password: "root",
		DB:       0,
	})

	return rdb
}
