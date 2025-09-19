package cart_service

import (
	"context"
	"fmt"

	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
)

type Repository interface {
	AddProduct(ctx context.Context, userId, productId string) (string, error)
	CartByUserId(ctx context.Context, userId string) (models.Cart, error)
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) AddProduct(ctx context.Context, userId, productId string) (string, error) {
	const op = "services.cart.AddProduct"

	cartId, err := s.repository.AddProduct(ctx, userId, productId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return cartId, nil
}

func (s *Service) CartByUserId(ctx context.Context, userId string) (models.Cart, error) {
	const op = "services.cart.CartByUserId"

	cart, err := s.repository.CartByUserId(ctx, userId)
	if err != nil {
		return models.Cart{}, fmt.Errorf("%s: %w", op, err)
	}

	for _, product := range cart.Products {
		cart.Price += product.Price
	}

	return cart, nil
}
