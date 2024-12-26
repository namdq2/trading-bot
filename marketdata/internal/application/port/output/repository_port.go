package output

import (
	"context"

	"marketdata/internal/domain/entity"
)

type OrderBookRepositoryPort interface {
	Store(ctx context.Context, orderbook *entity.OrderBook) error
	Get(ctx context.Context, exchangeID, symbol string) (*entity.OrderBook, error)
	GetAll(ctx context.Context) ([]*entity.OrderBook, error)
}

type TradeRepositoryPort interface {
	StoreTrade(ctx context.Context, trade *entity.Trade) error
	GetTradesBySymbol(ctx context.Context, symbol string, limit int) ([]*entity.Trade, error)
}

type EventPublisherPort interface {
	PublishOrderBookUpdate(ctx context.Context, orderbook *entity.OrderBook) error
	PublishTrade(ctx context.Context, trade *entity.Trade) error
}
