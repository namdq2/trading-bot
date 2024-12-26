package dto

import (
	"time"
)

type PriceLevelDTO struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

type OrderBookDTO struct {
	ExchangeID string          `json:"exchange_id"`
	Symbol     string          `json:"symbol"`
	Bids       []PriceLevelDTO `json:"bids"`
	Asks       []PriceLevelDTO `json:"asks"`
	Timestamp  time.Time       `json:"timestamp"`
}
