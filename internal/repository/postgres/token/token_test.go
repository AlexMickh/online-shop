package token_repository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/AlexMickh/coledzh-shop-backend/internal/consts"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgres_SaveToken(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx       context.Context
		userId    string
		token     string
		tokenType string
	}

	pool := initStorage()
	defer pool.Close()

	userId := uuid.NewString()
	_, _ = pool.Exec(context.Background(), "INSERT INTO users (id, email) VALUES ($1, $2)", userId, "sas-test223@gmail.com")

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
				ctx:       context.Background(),
				userId:    userId,
				token:     uuid.NewString(),
				tokenType: consts.TokenTypeEmailVerify,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				db: tt.fields.db,
			}
			if err := p.SaveToken(tt.args.ctx, tt.args.userId, tt.args.token, tt.args.tokenType); err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Postgres.SaveToken() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
	for _, tt := range tests {
		_, _ = pool.Exec(tt.args.ctx, "DELETE FROM tokens WHERE token = $1", tt.args.token)
	}

	_, _ = pool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userId)
}

func initStorage() *pgxpool.Pool {
	// TODO: change to config vars
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d",
		"postgres",
		"root",
		"localhost",
		5499,
		"store",
		3,
		5,
	)

	pool, _ := pgxpool.New(context.Background(), connString)

	return pool
}
