package category_service

import (
	"context"
	"fmt"
	"slices"
	"sort"

	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	SaveCategory(ctx context.Context, id string, name string) error
	AllCategories(ctx context.Context) ([]models.Category, error)
}

type Cash interface {
	SaveCategory(ctx context.Context, categoty models.Category) error
	AllCategories(ctx context.Context) ([]models.Category, error)
}

type Service struct {
	repository Repository
	cash       Cash
}

func New(repository Repository, cash Cash) *Service {
	return &Service{
		repository: repository,
		cash:       cash,
	}
}

func (s *Service) CreateCategory(ctx context.Context, name string) (string, error) {
	const op = "services.category.CreateCategory"

	id := uuid.NewString()
	err := s.repository.SaveCategory(ctx, id, name)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.cash.SaveCategory(ctx, models.Category{ID: id, Name: name})
	if err != nil {
		return id, fmt.Errorf("%s: %w", op, errs.ErrFailedToCash)
	}

	return id, nil
}

func (s *Service) AllCategories(ctx context.Context) ([]models.Category, error) {
	const op = "services.category.AllCategories"

	categories, err := s.cash.AllCategories(ctx)
	if err == nil {
		sortCategories(categories)
		return categories, nil
	}

	categories, err = s.repository.AllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	sortCategories(categories)

	return categories, nil
}

func sortCategories(arr []models.Category) {
	slices.SortFunc(arr, func(a models.Category, b models.Category) int {
		arr := []string{a.Name, b.Name}
		sort.Strings(arr)
		if arr[0] == a.Name {
			return -1
		}
		return 1
	})
}
