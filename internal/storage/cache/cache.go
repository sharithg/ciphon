package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Rdb *redis.Client
}

func New(rdb *redis.Client) *Cache {
	return &Cache{
		Rdb: rdb,
	}
}

func (c *Cache) StoreGithubCode(ctx context.Context) error {
	err := c.Rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetGithubCode(ctx context.Context, code string) (string, error) {
	val, err := c.Rdb.Get(ctx, "key").Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
