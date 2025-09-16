package user_service

import (
	"context"
	"fmt"

	"github.com/AlexMickh/coledzh-shop-backend/internal/consts"
	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
)

type SessionRepository interface {
	SessionById(ctx context.Context, sessionId string) (models.User, error)
}

type Service struct {
	sessionRepository SessionRepository
}

func New(sessionRepository SessionRepository) *Service {
	return &Service{
		sessionRepository: sessionRepository,
	}
}

func (s *Service) ValidateAdminSession(ctx context.Context, sessionId string) error {
	const op = "services.user.ValidateAdminSession"

	user, err := s.sessionRepository.SessionById(ctx, sessionId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if user.Role != consts.RoleAdmin {
		return fmt.Errorf("%s: %w", op, errs.ErrNotAdmin)
	}

	return nil
}
