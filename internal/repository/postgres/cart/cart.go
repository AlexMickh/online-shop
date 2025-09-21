package cart_repository

import (
	"context"
	"fmt"

	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
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

func (p *Postgres) AddProduct(ctx context.Context, userId, productId string) (string, error) {
	const op = "repository.postgres.cart.CartByUserId"

	var cartId string
	query := "INSERT INTO cart_items (user_id, product_id) VALUES ($1, $2) RETURNING id"
	err := p.db.QueryRow(ctx, query, userId, productId).Scan(&cartId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return cartId, nil
}

func (p *Postgres) CartByUserId(ctx context.Context, userId string) (models.Cart, error) {
	const op = "repository.postgres.cart.CartByUserId"

	var cart models.Cart
	cart.Products = make([]models.ProductCard, 0)
	query := `SELECT c.id, p.id, p.name, p.price, p.image_url
			  FROM cart_items c
			  JOIN products p
			  ON c.product_id = p.id
			  AND c.user_id = $1
			  GROUP BY c.id, p.id`
	rows, err := p.db.Query(ctx, query, userId)
	if err != nil {
		return models.Cart{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var product models.ProductCard
		err = rows.Scan(
			&cart.ID,
			&product.ID,
			&product.Name,
			&product.Price,
			&product.ImageUrl,
		)
		if err != nil {
			return models.Cart{}, fmt.Errorf("%s: %w", op, err)
		}
		cart.Products = append(cart.Products, product)
	}

	if rows.Err() != nil {
		return models.Cart{}, fmt.Errorf("%s: %w", op, err)
	}

	cart.UserId = userId

	return cart, nil
}

func (p *Postgres) DeleteCartByUserId(ctx context.Context, userId string) error {
	const op = "repository.postgres.cart.DeleteCartByUserId"

	query := "DELETE FROM cart_items WHERE user_id = $1"
	_, err := p.db.Exec(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
