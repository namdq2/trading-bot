package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"marketdata/internal/domain/entity"
)

type OrderBookRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewOrderBookRepository(client *redis.Client, ttl time.Duration) *OrderBookRepository {
	return &OrderBookRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *OrderBookRepository) Store(ctx context.Context, orderbook *entity.OrderBook) error {
	key := r.makeKey(orderbook.ExchangeID(), orderbook.Symbol())

	data, err := json.Marshal(orderbook)
	if err != nil {
		return fmt.Errorf("failed to marshal orderbook: %w", err)
	}

	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		return fmt.Errorf("failed to store orderbook in redis: %w", err)
	}

	return nil
}

func (r *OrderBookRepository) Get(ctx context.Context, exchangeID, symbol string) (*entity.OrderBook, error) {
	key := r.makeKey(exchangeID, symbol)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get orderbook from redis: %w", err)
	}

	var orderbook entity.OrderBook
	if err := json.Unmarshal(data, &orderbook); err != nil {
		return nil, fmt.Errorf("failed to unmarshal orderbook: %w", err)
	}

	return &orderbook, nil
}

func (r *OrderBookRepository) makeKey(exchangeID, symbol string) string {
	return fmt.Sprintf("orderbook:%s:%s", exchangeID, symbol)
}
