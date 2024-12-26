package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/yourusername/marketdata/internal/config"
	"github.com/yourusername/marketdata/internal/infrastructure/exchange"
	"github.com/yourusername/marketdata/internal/models"
)

type MarketDataService struct {
	logger     *zap.Logger
	config     *config.Config
	exchanges  map[string]exchange.Client
	orderbooks map[string]*models.OrderBook
	mu         sync.RWMutex
	done       chan struct{}
}

func NewMarketDataService(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*MarketDataService, error) {
	svc := &MarketDataService{
		logger:     logger,
		config:     cfg,
		exchanges:  make(map[string]exchange.Client),
		orderbooks: make(map[string]*models.OrderBook),
		done:       make(chan struct{}),
	}

	if err := svc.initializeExchanges(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize exchanges: %w", err)
	}

	return svc, nil
}

func (s *MarketDataService) Start(ctx context.Context) error {
	s.logger.Info("starting market data service")

	// Start orderbook subscriptions for each exchange and symbol
	for _, ex := range s.exchanges {
		for _, symbol := range s.config.Symbols {
			if err := s.subscribeOrderBook(ctx, ex, symbol); err != nil {
				return fmt.Errorf("failed to subscribe to orderbook: %w", err)
			}
		}
	}

	// Start maintenance routines
	go s.maintenance(ctx)

	return nil
}

func (s *MarketDataService) Stop() error {
	s.logger.Info("stopping market data service")
	close(s.done)

	// Close exchange connections
	for name, ex := range s.exchanges {
		if err := ex.Close(); err != nil {
			s.logger.Error("failed to close exchange connection",
				zap.String("exchange", name),
				zap.Error(err))
		}
	}

	return nil
}

func (s *MarketDataService) GetOrderBook(exchange, symbol string) (*models.OrderBook, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := fmt.Sprintf("%s-%s", exchange, symbol)
	ob, exists := s.orderbooks[key]
	if !exists {
		return nil, fmt.Errorf("orderbook not found for %s", key)
	}

	return ob.Copy(), nil
}

func (s *MarketDataService) subscribeOrderBook(ctx context.Context, ex exchange.Client, symbol string) error {
	updates, err := ex.SubscribeOrderBook(ctx, symbol)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-s.done:
				return
			case ob := <-updates:
				s.updateOrderBook(ex.Name(), symbol, ob)
			}
		}
	}()

	return nil
}

func (s *MarketDataService) updateOrderBook(exchange, symbol string, ob *models.OrderBook) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s-%s", exchange, symbol)
	s.orderbooks[key] = ob
}

func (s *MarketDataService) maintenance(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.done:
			return
		case <-ticker.C:
			s.cleanStaleData()
		}
	}
}

func (s *MarketDataService) cleanStaleData() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	staleThreshold := 5 * time.Minute

	for key, ob := range s.orderbooks {
		if now.Sub(ob.Timestamp) > staleThreshold {
			s.logger.Warn("removing stale orderbook",
				zap.String("key", key),
				zap.Time("timestamp", ob.Timestamp))
			delete(s.orderbooks, key)
		}
	}
}
