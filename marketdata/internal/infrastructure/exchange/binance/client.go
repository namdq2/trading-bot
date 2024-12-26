package binance

import (
	"context"
	"fmt"
	"sync"

	"marketdata/internal/domain/entity"
	"marketdata/internal/infrastructure/exchange"
	"marketdata/internal/infrastructure/websocket"
)

type Client struct {
	*exchange.BaseExchange
	wsConn      *websocket.Conn
	subscribers map[string][]chan *entity.OrderBook
	mu          sync.RWMutex
}

func NewClient(cfg exchange.Config) *Client {
	return &Client{
		BaseExchange: exchange.NewBaseExchange(cfg),
		subscribers:  make(map[string][]chan *entity.OrderBook),
	}
}

func (c *Client) Connect(ctx context.Context) error {
	if err := c.BaseExchange.Connect(ctx); err != nil {
		return err
	}

	// Implement Binance-specific connection logic
	return nil
}

func (c *Client) Close() error {
	if err := c.BaseExchange.Close(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Close all subscriber channels
	for _, channels := range c.subscribers {
		for _, ch := range channels {
			close(ch)
		}
	}
	c.subscribers = make(map[string][]chan *entity.OrderBook)

	return nil
}

func (c *Client) GetOrderBook(ctx context.Context, symbol string) (*entity.OrderBook, error) {
	// Implement Binance-specific orderbook fetching
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) SubscribeOrderBook(ctx context.Context, symbol string) (<-chan *entity.OrderBook, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := make(chan *entity.OrderBook, 100)
	c.subscribers[symbol] = append(c.subscribers[symbol], ch)

	// Implement Binance-specific subscription logic
	go c.handleOrderBookUpdates(ctx, symbol, ch)

	return ch, nil
}

func (c *Client) handleOrderBookUpdates(ctx context.Context, symbol string, ch chan<- *entity.OrderBook) {
	// Implement websocket message handling
}
