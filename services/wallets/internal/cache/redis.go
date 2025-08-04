// Package cache provides Redis caching functionality for transaction history.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/config"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/model"
	"github.com/go-redis/redis/v8"
)

// RedisClient interface defines the Redis operations for transaction caching
type RedisClient interface {
	GetTransactionHistory(ctx context.Context, userID string) ([]model.Transaction, error)
	SaveTransactionHistory(ctx context.Context, userID string, transactions []model.Transaction) error
	DeleteTransactionHistory(ctx context.Context, userID string) error
	Close() error
}

// redisClient implements RedisClient interface
type redisClient struct {
	client *redis.Client
	ttl    time.Duration
}

var (
	redisInstance RedisClient
	redisOnce     sync.Once
)

// ResetRedisClient resets the singleton instance for testing
func ResetRedisClient() {
	redisOnce = sync.Once{}
	redisInstance = nil
}

// NewRedisClient creates a new Redis client instance using singleton pattern
func NewRedisClient() RedisClient {
	redisOnce.Do(func() {
		globalConfig := config.GetGlobalConfig()
		redisConfig := globalConfig.Redis

		rdb := redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
			Password:     redisConfig.Password,
			DB:           redisConfig.DB,
			MaxRetries:   redisConfig.MaxRetries,
			PoolSize:     redisConfig.PoolSize,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		})

		redisInstance = &redisClient{
			client: rdb,
			ttl:    24 * time.Hour, // Cache for 24 hours
		}
	})
	return redisInstance
}

// generateKey creates a unique Redis key for user transaction history
func (r *redisClient) generateKey(userID string) string {
	return fmt.Sprintf("wallet:transactions:%s", userID)
}

// GetTransactionHistory retrieves cached transaction history for a user
func (r *redisClient) GetTransactionHistory(ctx context.Context, userID string) ([]model.Transaction, error) {
	key := r.generateKey(userID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Cache miss - return empty slice
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get transaction history from cache: %w", err)
	}

	var transactions []model.Transaction
	err = json.Unmarshal([]byte(val), &transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction history: %w", err)
	}

	return transactions, nil
}

// SaveTransactionHistory caches transaction history for a user
func (r *redisClient) SaveTransactionHistory(ctx context.Context, userID string, transactions []model.Transaction) error {
	key := r.generateKey(userID)
	data, err := json.Marshal(transactions)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction history: %w", err)
	}

	err = r.client.Set(ctx, key, data, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save transaction history to cache: %w", err)
	}

	return nil
}

// DeleteTransactionHistory removes cached transaction history for a user
func (r *redisClient) DeleteTransactionHistory(ctx context.Context, userID string) error {
	key := r.generateKey(userID)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete transaction history from cache: %w", err)
	}

	return nil
}

// Close closes the Redis client connection
func (r *redisClient) Close() error {
	return r.client.Close()
}
