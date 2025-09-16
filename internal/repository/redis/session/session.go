package session_cash

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

func (c *Cash) SaveSession(ctx context.Context, id string, user models.User) error {
	const op = "repository.redis.session.SaveSession"

	err := c.rdb.HSet(ctx, id, user).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = c.rdb.Expire(ctx, id, c.expire).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Cash) SessionById(ctx context.Context, sessionId string) (models.User, error) {
	const op = "repository.redis.session.SessionById"

	var user models.User
	err := c.rdb.HGetAll(ctx, sessionId).Scan(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
