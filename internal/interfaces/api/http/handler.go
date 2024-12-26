package http

import (
	"encoding/json"
	"net/http"

	"marketdata/internal/application/port/input"
)

type MarketDataHandler struct {
	marketDataUseCase input.MarketDataUseCase
}

func NewMarketDataHandler(useCase input.MarketDataUseCase) *MarketDataHandler {
	return &MarketDataHandler{
		marketDataUseCase: useCase,
	}
}

func (h *MarketDataHandler) GetOrderBook(w http.ResponseWriter, r *http.Request) {
	exchangeID := r.URL.Query().Get("exchange")
	symbol := r.URL.Query().Get("symbol")

	if exchangeID == "" || symbol == "" {
		http.Error(w, "missing required parameters", http.StatusBadRequest)
		return
	}

	orderbook, err := h.marketDataUseCase.GetOrderBook(r.Context(), exchangeID, symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if orderbook == nil {
		http.Error(w, "orderbook not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderbook)
}
