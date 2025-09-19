package product_service

import (
	"context"
	"fmt"

	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	SaveProduct(
		ctx context.Context,
		productId string,
		name string,
		description string,
		price float32,
		imageUrl string,
		categoryIds []string,
	) error
	ProductsByCategoryId(ctx context.Context, categoryId string, page int) ([]models.ProductCard, error)
	AllProducts(ctx context.Context, page int) ([]models.ProductCard, error)
	ProductById(ctx context.Context, productId string) (models.Product, error)
}

type S3 interface {
	SaveImage(ctx context.Context, id string, image []byte) (string, error)
}

type Service struct {
	repository Repository
	s3         S3
}

func New(repository Repository, s3 S3) *Service {
	return &Service{
		repository: repository,
		s3:         s3,
	}
}

func (s *Service) CreateProduct(
	ctx context.Context,
	categoryIds []string,
	name string,
	description string,
	price float32,
	image []byte,
) (string, error) {
	const op = "services.product.CreateProduct"

	productId := uuid.NewString()

	imageUrl, err := s.s3.SaveImage(ctx, productId, image)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.repository.SaveProduct(
		ctx,
		productId,
		name,
		description,
		price,
		imageUrl,
		categoryIds,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return productId, nil
}

func (s *Service) ProductsCard(ctx context.Context, categoryId string, page int) ([]models.ProductCard, error) {
	const op = "services.product.CreateProduct"

	if categoryId != "" {
		products, err := s.repository.ProductsByCategoryId(ctx, categoryId, page)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return products, nil
	}

	products, err := s.repository.AllProducts(ctx, page)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (s *Service) ProductById(ctx context.Context, productId string) (models.Product, error) {
	const op = "services.product.ProductById"

	product, err := s.repository.ProductById(ctx, productId)
	if err != nil {
		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}

	return product, nil
}
