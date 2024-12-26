package dto

import (
	"time"
)

type TradeDTO struct {
	ID         string    `json:"id"`
	ExchangeID string    `json:"exchange_id"`
	Symbol     string    `json:"symbol"`
	Price      float64   `json:"price"`
	Volume     float64   `json:"volume"`
	TradeType  string    `json:"trade_type"`
	Timestamp  time.Time `json:"timestamp"`
}
