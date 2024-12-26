package valueobject

import (
	"fmt"
	"math"
)

// Price represents a price value object
type Price struct {
	value    float64
	currency string
}

func NewPrice(value float64, currency string) (*Price, error) {
	if value < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}
	if currency == "" {
		return nil, fmt.Errorf("currency cannot be empty")
	}

	return &Price{
		value:    value,
		currency: currency,
	}, nil
}

func (p Price) Value() float64 {
	return p.value
}

func (p Price) Currency() string {
	return p.currency
}

func (p Price) Equals(other Price) bool {
	return p.currency == other.currency && math.Abs(p.value-other.value) < 0.00000001
}
