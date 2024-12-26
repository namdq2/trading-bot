package exchange

import (
	"context"

	"github.com/yourusername/marketdata/internal/models"
)

// Client defines the interface that all exchange clients must implement
type Client interface {
	// Connect establishes connection to the exchange
	Connect(ctx context.Context) error

	// Close closes the connection to the exchange
	Close() error

	// SubscribeOrderBook subscribes to orderbook updates for a symbol
	SubscribeOrderBook(ctx context.Context, symbol string) (<-chan *models.OrderBook, error)

	// GetOrderBook gets a snapshot of the current orderbook
	GetOrderBook(ctx context.Context, symbol string) (*models.OrderBook, error)

	// Name returns the exchange name
	Name() string
}

// Config contains common exchange configuration
type Config struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	WSEndpoint string
}
