package user_repository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgres_SaveUser(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx      context.Context
		login    string
		email    string
		password string
	}

	pool := initStorage()
	defer pool.Close()

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
				ctx:      context.Background(),
				login:    "sas",
				email:    "sas2@gmail.com",
				password: "sas123",
			},
			wantErr: nil,
		},
		{
			name: "case user already exists",
			fields: fields{
				db: pool,
			},
			args: args{
				ctx:      context.Background(),
				login:    "sas",
				email:    "sas2@gmail.com",
				password: "sas123",
			},
			wantErr: errs.ErrUserAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				db: tt.fields.db,
			}
			got, err := p.SaveUser(tt.args.ctx, tt.args.login, tt.args.email, tt.args.password)
			if err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Postgres.SaveUser() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if got == "" && tt.name != "case user already exists" {
				t.Errorf("Postgres.SaveUser() = id is empty")
			}
		})
	}
	for _, tt := range tests {
		_, _ = pool.Exec(tt.args.ctx, "DELETE FROM users WHERE email = $1", tt.args.email)
	}
}

func TestPostgres_SaveAdmin(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx      context.Context
		login    string
		email    string
		password string
	}

	pool := initStorage()
	defer pool.Close()

	email := gofakeit.Email()

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
				ctx:      context.Background(),
				login:    "admin2",
				email:    email,
				password: "admin2",
			},
			wantErr: nil,
		},
		{
			name: "case user already exists",
			fields: fields{
				db: pool,
			},
			args: args{
				ctx:      context.Background(),
				login:    "admin2",
				email:    email,
				password: "admin2",
			},
			wantErr: errs.ErrUserAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				db: tt.fields.db,
			}
			if err := p.SaveAdmin(tt.args.ctx, tt.args.login, tt.args.email, tt.args.password); err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Postgres.SaveAdmin() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
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
