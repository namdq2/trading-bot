package entity

import (
	"time"

	"marketdata/internal/domain/valueobject"
)

type TradeType string

const (
	TradeBuy  TradeType = "BUY"
	TradeSell TradeType = "SELL"
)

type Trade struct {
	id         string
	exchangeID string
	symbol     string
	price      valueobject.Price
	volume     valueobject.Volume
	tradeType  TradeType
	timestamp  time.Time
}

func NewTrade(
	id string,
	exchangeID string,
	symbol string,
	price valueobject.Price,
	volume valueobject.Volume,
	tradeType TradeType,
	timestamp time.Time,
) *Trade {
	return &Trade{
		id:         id,
		exchangeID: exchangeID,
		symbol:     symbol,
		price:      price,
		volume:     volume,
		tradeType:  tradeType,
		timestamp:  timestamp,
	}
}

func (t *Trade) ID() string {
	return t.id
}

func (t *Trade) ExchangeID() string {
	return t.exchangeID
}

func (t *Trade) Symbol() string {
	return t.symbol
}

func (t *Trade) Price() valueobject.Price {
	return t.price
}

func (t *Trade) Volume() valueobject.Volume {
	return t.volume
}

func (t *Trade) Type() TradeType {
	return t.tradeType
}

func (t *Trade) Timestamp() time.Time {
	return t.timestamp
}
