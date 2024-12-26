package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	orderBookUpdates    *prometheus.CounterVec
	orderBookLatency    *prometheus.HistogramVec
	tradeUpdates        *prometheus.CounterVec
	exchangeErrors      *prometheus.CounterVec
	activeSubscriptions *prometheus.GaugeVec
}

func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		orderBookUpdates: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "orderbook_updates_total",
				Help:      "Total number of orderbook updates received",
			},
			[]string{"exchange", "symbol"},
		),
		orderBookLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "orderbook_latency_seconds",
				Help:      "Latency of orderbook updates",
				Buckets:   prometheus.ExponentialBuckets(0.001, 2, 10),
			},
			[]string{"exchange", "symbol"},
		),
		tradeUpdates: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "trade_updates_total",
				Help:      "Total number of trade updates received",
			},
			[]string{"exchange", "symbol"},
		),
		exchangeErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "exchange_errors_total",
				Help:      "Total number of exchange errors",
			},
			[]string{"exchange", "type"},
		),
		activeSubscriptions: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "active_subscriptions",
				Help:      "Number of active subscriptions",
			},
			[]string{"exchange", "symbol"},
		),
	}
}

func (m *Metrics) RecordOrderBookUpdate(exchange, symbol string) {
	m.orderBookUpdates.WithLabelValues(exchange, symbol).Inc()
}

func (m *Metrics) RecordOrderBookLatency(exchange, symbol string, latency float64) {
	m.orderBookLatency.WithLabelValues(exchange, symbol).Observe(latency)
}

func (m *Metrics) RecordTradeUpdate(exchange, symbol string) {
	m.tradeUpdates.WithLabelValues(exchange, symbol).Inc()
}

func (m *Metrics) RecordExchangeError(exchange, errorType string) {
	m.exchangeErrors.WithLabelValues(exchange, errorType).Inc()
}

func (m *Metrics) SetActiveSubscriptions(exchange, symbol string, count float64) {
	m.activeSubscriptions.WithLabelValues(exchange, symbol).Set(count)
}
