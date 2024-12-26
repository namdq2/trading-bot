package input

import (
	"context"

	"marketdata/internal/application/dto"
)

type MarketDataUseCase interface {
	// GetOrderBook retrieves the current orderbook for a given exchange and symbol
	GetOrderBook(ctx context.Context, exchangeID, symbol string) (*dto.OrderBookDTO, error)

	// SubscribeOrderBook subscribes to orderbook updates for a given exchange and symbol
	SubscribeOrderBook(ctx context.Context, exchangeID, symbol string) (<-chan *dto.OrderBookDTO, error)

	// GetTrades retrieves recent trades for a given exchange and symbol
	GetTrades(ctx context.Context, exchangeID, symbol string, limit int) ([]*dto.TradeDTO, error)

	// ProcessOrderBookUpdate processes an orderbook update from an exchange
	ProcessOrderBookUpdate(ctx context.Context, update *dto.OrderBookDTO) error
}
