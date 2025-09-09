package category_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	SaveCategory(ctx context.Context, id string, name string) error
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateCategory(ctx context.Context, name string) (string, error) {
	const op = "services.category.CreateCategory"

	id := uuid.NewString()
	err := s.repository.SaveCategory(ctx, id, name)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
