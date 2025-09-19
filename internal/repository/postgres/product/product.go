package product_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

const pageSize = 10

func New(db *pgxpool.Pool) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) SaveProduct(
	ctx context.Context,
	productId string,
	name string,
	description string,
	price float32,
	imageUrl string,
	categoryIds []string,
) error {
	const op = "repository.postgres.product.SaveProduct"

	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	query := `INSERT INTO products
			  (id, name, description, price, image_url)
			  VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(ctx, query, productId, name, description, price, imageUrl)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	builder := new(strings.Builder)
	argsCounter := 1

	_, err = builder.WriteString("INSERT INTO products_categories (category_id, product_id) VALUES ")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	l := len(categoryIds)
	args := make([]any, 0, l*2)

	for i, categoryId := range categoryIds {
		_, err = builder.WriteString(fmt.Sprintf("($%d, $%d) ", argsCounter, argsCounter+1))
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if i < l-1 {
			_, err = builder.WriteString(", ")
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

		argsCounter += 2
		args = append(args, categoryId)
		args = append(args, productId)
	}

	_, err = tx.Exec(ctx, builder.String(), args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// TODO: add tests
func (p *Postgres) ProductsByCategoryId(ctx context.Context, categoryId string, page int) ([]models.ProductCard, error) {
	const op = "repository.postgres.product.ProductsByCategoryId"

	query := `SELECT p.id, p.name, p.price, p.image_url
			  FROM products p
			  JOIN products_categories pc
			  ON p.id = pc.product_id
			  AND pc.category_id = $1
			  ORDER BY p.price
			  OFFSET $2
			  LIMIT $3`
	products := make([]models.ProductCard, 0, pageSize)
	rows, err := p.db.Query(ctx, query, categoryId, page*pageSize, page*pageSize+pageSize)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var product models.ProductCard

		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.ImageUrl,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (p *Postgres) AllProducts(ctx context.Context, page int) ([]models.ProductCard, error) {
	const op = "repository.postgres.product.AllProducts"

	query := `SELECT p.id, p.name, p.price, p.image_url
			  FROM products p
			  JOIN products_categories pc
			  ON p.id = pc.product_id
			  ORDER BY p.price
			  OFFSET $1
			  LIMIT $2`
	products := make([]models.ProductCard, 0, pageSize)
	rows, err := p.db.Query(ctx, query, page*pageSize, page*pageSize+pageSize)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var product models.ProductCard

		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.ImageUrl,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (p *Postgres) ProductById(ctx context.Context, productId string) (models.Product, error) {
	const op = "repository.postgres.product.AllProducts"

	var product models.Product
	var categoryIds string
	var categoryNames string

	query := `SELECT p.id, p.name, p.description, p.price, p.image_url, 
	              string_agg(c.id::text, ' ') AS category_idss, string_agg(c.name, ' ') AS category_names
			  FROM products AS p
			  JOIN products_categories AS pc
			  ON p.id = pc.product_id
			  JOIN categories AS c
			  ON pc.category_id = c.id
			  AND p.id = $1
			  GROUP BY p.id`
	err := p.db.QueryRow(ctx, query, productId).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.ImageUrl,
		&categoryIds,
		&categoryNames,
	)
	if err != nil {
		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}

	categoryNamesArr := strings.Split(categoryNames, " ")
	product.Categories = make([]models.Category, 0)

	for i, categoryId := range strings.Split(categoryIds, " ") {
		category := models.Category{
			ID:   categoryId,
			Name: categoryNamesArr[i],
		}
		product.Categories = append(product.Categories, category)
	}

	return product, nil
}
