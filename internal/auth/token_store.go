package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RefreshTokenStore interface {
	CreateRefreshToken(ctx context.Context, userID string, ttl time.Duration) (string, error)
	VerifyRefreshToken(ctx context.Context, token string) (string, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

type TokenStore struct {
	Redis *redis.Client
}

func NewTokenStore(addr string) *TokenStore {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &TokenStore{Redis: rdb}
}

func (ts *TokenStore) CreateRefreshToken(ctx context.Context, userID string, ttl time.Duration) (string, error) {
	token := uuid.New().String()
	err := ts.Redis.Set(ctx, token, userID, ttl).Err()
	return token, err
}

func (ts *TokenStore) VerifyRefreshToken(ctx context.Context, token string) (string, error) {
	return ts.Redis.Get(ctx, token).Result()
}

func (ts *TokenStore) DeleteRefreshToken(ctx context.Context, token string) error {
	return ts.Redis.Del(ctx, token).Err()
}
