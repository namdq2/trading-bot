package entity

import (
	"time"

	"marketdata/internal/domain/valueobject"
)

type PriceLevel struct {
	Price    valueobject.Price
	Quantity valueobject.Volume
}

type OrderBook struct {
	exchangeID string
	symbol     string
	bids       []PriceLevel
	asks       []PriceLevel
	timestamp  time.Time
}

func NewOrderBook(exchangeID, symbol string, timestamp time.Time) *OrderBook {
	return &OrderBook{
		exchangeID: exchangeID,
		symbol:     symbol,
		bids:       make([]PriceLevel, 0),
		asks:       make([]PriceLevel, 0),
		timestamp:  timestamp,
	}
}

func (ob *OrderBook) UpdateBids(bids []PriceLevel) {
	ob.bids = bids
}

func (ob *OrderBook) UpdateAsks(asks []PriceLevel) {
	ob.asks = asks
}

func (ob *OrderBook) BestBid() (PriceLevel, bool) {
	if len(ob.bids) == 0 {
		return PriceLevel{}, false
	}
	return ob.bids[0], true
}

func (ob *OrderBook) BestAsk() (PriceLevel, bool) {
	if len(ob.asks) == 0 {
		return PriceLevel{}, false
	}
	return ob.asks[0], true
}
