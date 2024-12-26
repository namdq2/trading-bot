package service

import (
	"marketdata/internal/domain/entity"
)

type OrderBookDomainService interface {
	ValidateOrderBook(orderbook *entity.OrderBook) error
	CalculateSpread(orderbook *entity.OrderBook) (float64, error)
	DetectArbitrageOpportunity(books []*entity.OrderBook) ([]*ArbitrageOpportunity, error)
}

type TradeDomainService interface {
	ValidateTrade(trade *entity.Trade) error
	CalculateTradeValue(trade *entity.Trade) (float64, error)
}

type ArbitrageOpportunity struct {
	BuyExchange  string
	SellExchange string
	Symbol       string
	SpreadPct    float64
	MaxVolume    float64
}
