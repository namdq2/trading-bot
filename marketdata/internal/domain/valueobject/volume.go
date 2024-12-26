package valueobject

import (
	"fmt"
	"math"
)

// Volume represents a quantity/amount value object
type Volume struct {
	value float64
	asset string
}

func NewVolume(value float64, asset string) (*Volume, error) {
	if value < 0 {
		return nil, fmt.Errorf("volume cannot be negative")
	}
	if asset == "" {
		return nil, fmt.Errorf("asset cannot be empty")
	}

	return &Volume{
		value: value,
		asset: asset,
	}, nil
}

func (v Volume) Value() float64 {
	return v.value
}

func (v Volume) Asset() string {
	return v.asset
}

func (v Volume) Add(other Volume) (*Volume, error) {
	if v.asset != other.asset {
		return nil, fmt.Errorf("cannot add volumes of different assets")
	}
	return NewVolume(v.value+other.value, v.asset)
}

func (v Volume) Subtract(other Volume) (*Volume, error) {
	if v.asset != other.asset {
		return nil, fmt.Errorf("cannot subtract volumes of different assets")
	}
	return NewVolume(v.value-other.value, v.asset)
}

func (v Volume) Equals(other Volume) bool {
	return v.asset == other.asset && math.Abs(v.value-other.value) < 0.00000001
}
