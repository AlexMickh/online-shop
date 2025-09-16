package category_cash

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	"github.com/redis/go-redis/v9"
)

type Cash struct {
	rdb    *redis.Client
	expire time.Duration
}

func New(rdb *redis.Client, expire time.Duration) *Cash {
	return &Cash{
		rdb:    rdb,
		expire: expire,
	}
}

func (c *Cash) SaveCategory(ctx context.Context, categoty models.Category) error {
	const op = "repository.redis.category.SaveCategory"

	err := c.rdb.HSet(ctx, genKey(categoty.ID), categoty).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Cash) AllCategories(ctx context.Context) ([]models.Category, error) {
	const op = "repository.redis.category.AllCategories"

	var (
		err    error
		cursor uint64
		keys   []string
	)
	for {
		keys, cursor, err = c.rdb.Scan(ctx, cursor, "category:*", 10).Result()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if cursor == 0 {
			break
		}
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("%s: nothing found", op)
	}

	categories := make([]models.Category, 0, len(keys))
	for _, key := range keys {
		var category models.Category
		err = c.rdb.HGetAll(ctx, key).Scan(&category)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		category.ID = key
		categories = append(categories, category)
	}

	return categories, nil
}

func genKey(id string) string {
	return "category:" + id
}
