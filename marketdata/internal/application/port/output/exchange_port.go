package output

import (
	"context"

	"marketdata/internal/domain/entity"
)

// ExchangePort defines the interface that exchange adapters must implement
type ExchangePort interface {
	// Connect establishes connection to the exchange
	Connect(ctx context.Context) error

	// Close closes the connection to the exchange
	Close() error

	// GetOrderBook gets a snapshot of the current orderbook
	GetOrderBook(ctx context.Context, symbol string) (*entity.OrderBook, error)

	// SubscribeOrderBook subscribes to orderbook updates for a symbol
	SubscribeOrderBook(ctx context.Context, symbol string) (<-chan *entity.OrderBook, error)

	// GetName returns the exchange name
	GetName() string
}
