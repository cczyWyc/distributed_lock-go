package demo

import (
	"context"
	_ "embed"
	"errors"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"time"
)

var (
	//go:embed unlock.lua
	luaUnlock string
)

type Client struct {
	client redis.Cmdable
}

type Lock struct {
	client redis.Cmdable
	key    string
	value  string
}

// NewClient create a new client instance
func NewClient(c redis.Cmdable) *Client {
	return &Client{
		client: c,
	}
}

// NewLock create a new instance
func newLock(client redis.Cmdable, key string, value string) *Lock {
	return &Lock{
		client: client,
		key:    key,
		value:  value,
	}
}

// Lock add lock with key
func (c *Client) Lock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {
	value := uuid.New().String()
	res, err := c.client.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, errors.New("add lock failed")
	}

	return newLock(c.client, key, value), nil
}

// Unlock unlock a lock
func (l *Lock) Unlock(ctx context.Context) error {
	l.client.Eval(ctx, luaUnlock, []string{l.key}, l.value)

	return nil
}
