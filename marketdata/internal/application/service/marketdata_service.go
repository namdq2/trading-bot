package service

import (
	"context"
	"fmt"

	"marketdata/internal/application/dto"
	"marketdata/internal/domain/entity"
	"marketdata/internal/domain/repository"
)

type MarketDataService struct {
	orderbookRepo repository.OrderBookRepository
	tradeRepo     repository.TradeRepository
	publisher     EventPublisher
	logger        Logger
}

type EventPublisher interface {
	PublishOrderBookUpdate(ctx context.Context, orderbook *dto.OrderBookDTO) error
}

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

func NewMarketDataService(
	orderbookRepo repository.OrderBookRepository,
	tradeRepo repository.TradeRepository,
	publisher EventPublisher,
	logger Logger,
) *MarketDataService {
	return &MarketDataService{
		orderbookRepo: orderbookRepo,
		tradeRepo:     tradeRepo,
		publisher:     publisher,
		logger:        logger,
	}
}

func (s *MarketDataService) ProcessOrderBookUpdate(ctx context.Context, update *dto.OrderBookDTO) error {
	// Convert DTO to domain entity
	orderbook := entity.NewOrderBook(
		update.ExchangeID,
		update.Symbol,
		update.Timestamp,
	)

	// Store in repository
	if err := s.orderbookRepo.Store(ctx, orderbook); err != nil {
		return fmt.Errorf("failed to store orderbook: %w", err)
	}

	// Publish update
	if err := s.publisher.PublishOrderBookUpdate(ctx, update); err != nil {
		s.logger.Error("failed to publish orderbook update",
			"error", err,
			"exchange", update.ExchangeID,
			"symbol", update.Symbol,
		)
	}

	return nil
}
