package exchange

import (
	"context"
	"fmt"
	"sync"

	"marketdata/internal/application/port/output"
)

// BaseExchange provides common functionality for exchange implementations
type BaseExchange struct {
	name      string
	apiKey    string
	apiSecret string
	baseURL   string
	wsURL     string
	mu        sync.RWMutex
	connected bool
}

// Config contains common exchange configuration
type Config struct {
	Name      string
	APIKey    string
	APISecret string
	BaseURL   string
	WSURL     string
}

// NewBaseExchange creates a new BaseExchange
func NewBaseExchange(cfg Config) *BaseExchange {
	return &BaseExchange{
		name:      cfg.Name,
		apiKey:    cfg.APIKey,
		apiSecret: cfg.APISecret,
		baseURL:   cfg.BaseURL,
		wsURL:     cfg.WSURL,
	}
}

// GetName returns the exchange name
func (e *BaseExchange) GetName() string {
	return e.name
}

// IsConnected returns the connection status
func (e *BaseExchange) IsConnected() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.connected
}

// setConnected sets the connection status
func (e *BaseExchange) setConnected(status bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.connected = status
}

// Connect implements the base connection logic
func (e *BaseExchange) Connect(ctx context.Context) error {
	if e.IsConnected() {
		return fmt.Errorf("exchange %s is already connected", e.name)
	}
	e.setConnected(true)
	return nil
}

// Close implements the base close logic
func (e *BaseExchange) Close() error {
	if !e.IsConnected() {
		return fmt.Errorf("exchange %s is not connected", e.name)
	}
	e.setConnected(false)
	return nil
}

// Ensure BaseExchange implements ExchangePort
var _ output.ExchangePort = (*BaseExchange)(nil)
