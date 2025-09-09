package auth_service

import (
	"context"
	"fmt"

	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	SaveUser(ctx context.Context, login, email, password string) (string, error)
	SaveAdmin(ctx context.Context, login, email, password string) error
	UserByEmail(ctx context.Context, email string) (models.User, error)
	VerifyEmail(ctx context.Context, id string) error
}

type SessionStore interface {
	SaveSession(ctx context.Context, id string, user models.User) error
}

type Service struct {
	storage      Storage
	sessionStore SessionStore
}

func New(storage Storage, sessionStore SessionStore) *Service {
	return &Service{
		storage:      storage,
		sessionStore: sessionStore,
	}
}

func (s *Service) Register(ctx context.Context, login, email, password string) (string, error) {
	const op = "services.auth.Register"

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.storage.SaveUser(ctx, login, email, string(hashPassword))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Service) RegisterAdmin(ctx context.Context, login, email, password string) error {
	const op = "services.auth.Register"

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.storage.SaveAdmin(ctx, login, email, string(hashPassword))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	const op = "services.auth.Login"

	user, err := s.storage.UserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !user.IsEmailVerified {
		return "", fmt.Errorf("%s: %w", op, errs.ErrEmailNotVerify)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, errs.ErrUserNotFound)
	}

	sessionId := uuid.NewString()
	err = s.sessionStore.SaveSession(ctx, sessionId, user)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return sessionId, nil
}

// TODO: move to user service
func (s *Service) VerifyEmail(ctx context.Context, id string) error {
	const op = "services.auth.VerifyEmail"

	err := s.storage.VerifyEmail(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
