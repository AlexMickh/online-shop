package token_service

import (
	"context"
	"fmt"

	"github.com/AlexMickh/coledzh-shop-backend/internal/consts"
	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	"github.com/google/uuid"
)

type Storage interface {
	SaveToken(ctx context.Context, userId, token, tokenType string) error
	UserIdByToken(ctx context.Context, token, tokenType string) (string, error)
}

type UserService interface {
	VerifyEmail(ctx context.Context, id string) error
}

type Service struct {
	storage     Storage
	userService UserService
}

func New(storage Storage, userService UserService) *Service {
	return &Service{
		storage:     storage,
		userService: userService,
	}
}

func (s *Service) CreateToken(ctx context.Context, userId, tokenType string) (string, error) {
	const op = "services.token.CreateToken"

	var token string
	switch tokenType {
	case consts.TokenTypeEmailVerify:
		token = uuid.NewString()
	default:
		return "", fmt.Errorf("%s: %w", op, errs.ErrWrongTokenType)
	}

	err := s.storage.SaveToken(ctx, userId, token, tokenType)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	const op = "services.token.ValidateEmail"

	userId, err := s.storage.UserIdByToken(ctx, token, consts.TokenTypeEmailVerify)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.userService.VerifyEmail(ctx, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
