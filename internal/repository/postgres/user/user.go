package user_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/AlexMickh/coledzh-shop-backend/internal/consts"
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

func (p *Postgres) SaveUser(ctx context.Context, login, email, password string) (string, error) {
	const op = "repository.postgres.user.SaveUser"

	query := `INSERT INTO users 
			  (login, email, password)
			  VALUES ($1, $2, $3)
			  RETURNING id`
	var id string
	err := p.db.QueryRow(ctx, query, login, email, password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return "", fmt.Errorf("%s: %w", op, errs.ErrUserAlreadyExists)
			}
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (p *Postgres) SaveAdmin(ctx context.Context, login, email, password string) error {
	const op = "repository.postgres.user.SaveAdmin"

	query := `INSERT INTO users 
			  (login, email, password, role, is_email_verified)
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := p.db.Exec(ctx, query, login, email, password, consts.RoleAdmin, true)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, errs.ErrUserAlreadyExists)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *Postgres) UserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "repository.postgres.user.UserByEmail"

	query := "SELECT id, login, email, password, role, is_email_verified FROM users WHERE email = $1"
	var user models.User
	err := p.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsEmailVerified,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, errs.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (p *Postgres) VerifyEmail(ctx context.Context, id string) error {
	const op = "repository.postgres.user.VerifyEmail"

	query := "UPDATE users SET is_email_verified = $1 WHERE id = $2"
	_, err := p.db.Exec(ctx, query, true, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
