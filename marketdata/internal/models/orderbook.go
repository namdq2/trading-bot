package models

import (
	"time"
)

// PriceLevel represents a single price level in the orderbook
type PriceLevel struct {
	Price     float64
	Quantity  float64
	Timestamp time.Time
}

// OrderBook represents the full orderbook for a trading pair
type OrderBook struct {
	Exchange  string
	Symbol    string
	Bids      []PriceLevel
	Asks      []PriceLevel
	Timestamp time.Time
}

// Copy creates a deep copy of the orderbook
func (ob *OrderBook) Copy() *OrderBook {
	copy := &OrderBook{
		Exchange:  ob.Exchange,
		Symbol:    ob.Symbol,
		Timestamp: ob.Timestamp,
		Bids:      make([]PriceLevel, len(ob.Bids)),
		Asks:      make([]PriceLevel, len(ob.Asks)),
	}

	for i := range ob.Bids {
		copy.Bids[i] = ob.Bids[i]
	}
	for i := range ob.Asks {
		copy.Asks[i] = ob.Asks[i]
	}

	return copy
}

// BestBid returns the highest bid price and quantity
func (ob *OrderBook) BestBid() (price, quantity float64) {
	if len(ob.Bids) == 0 {
		return 0, 0
	}
	return ob.Bids[0].Price, ob.Bids[0].Quantity
}

// BestAsk returns the lowest ask price and quantity
func (ob *OrderBook) BestAsk() (price, quantity float64) {
	if len(ob.Asks) == 0 {
		return 0, 0
	}
	return ob.Asks[0].Price, ob.Asks[0].Quantity
}
