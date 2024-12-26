package service

import (
	"context"
	"fmt"

	"marketdata/internal/application/dto"
	"marketdata/internal/application/port/input"
	"marketdata/internal/application/port/output"
	"marketdata/internal/domain/entity"
)

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

type MarketDataService struct {
	orderbookRepo output.OrderBookRepositoryPort
	tradeRepo     output.TradeRepositoryPort
	exchangeMgr   output.ExchangePort
	publisher     output.EventPublisherPort
	logger        Logger
}

func NewMarketDataService(
	orderbookRepo output.OrderBookRepositoryPort,
	tradeRepo output.TradeRepositoryPort,
	exchangeMgr output.ExchangePort,
	publisher output.EventPublisherPort,
	logger Logger,
) *MarketDataService {
	return &MarketDataService{
		orderbookRepo: orderbookRepo,
		tradeRepo:     tradeRepo,
		exchangeMgr:   exchangeMgr,
		publisher:     publisher,
		logger:        logger,
	}
}

// GetOrderBook retrieves the current orderbook for a given exchange and symbol
func (s *MarketDataService) GetOrderBook(ctx context.Context, exchangeID, symbol string) (*dto.OrderBookDTO, error) {
	// Get orderbook from repository
	orderbook, err := s.orderbookRepo.Get(ctx, exchangeID, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get orderbook: %w", err)
	}

	if orderbook == nil {
		return nil, nil
	}

	// Convert domain entity to DTO
	return convertToOrderBookDTO(orderbook), nil
}

// SubscribeOrderBook subscribes to orderbook updates for a given exchange and symbol
func (s *MarketDataService) SubscribeOrderBook(ctx context.Context, exchangeID, symbol string) (<-chan *dto.OrderBookDTO, error) {
	// Subscribe to exchange updates
	updates, err := s.exchangeMgr.SubscribeOrderBook(ctx, exchangeID, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to orderbook: %w", err)
	}

	// Create channel for DTOs
	dtoChan := make(chan *dto.OrderBookDTO, 100)

	// Start goroutine to convert and forward updates
	go func() {
		defer close(dtoChan)
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-updates:
				if update == nil {
					return
				}
				dtoChan <- convertToOrderBookDTO(update)
			}
		}
	}()

	return dtoChan, nil
}

// GetTrades retrieves recent trades for a given exchange and symbol
func (s *MarketDataService) GetTrades(ctx context.Context, exchangeID, symbol string, limit int) ([]*dto.TradeDTO, error) {
	// Get trades from repository
	trades, err := s.tradeRepo.GetTradesBySymbol(ctx, symbol, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades: %w", err)
	}

	// Convert domain entities to DTOs
	tradeDTOs := make([]*dto.TradeDTO, len(trades))
	for i, trade := range trades {
		tradeDTOs[i] = convertToTradeDTO(trade)
	}

	return tradeDTOs, nil
}

// ProcessOrderBookUpdate processes an orderbook update from an exchange
func (s *MarketDataService) ProcessOrderBookUpdate(ctx context.Context, update *dto.OrderBookDTO) error {
	// Convert DTO to domain entity
	orderbook := convertToOrderBookEntity(update)

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

// Helper functions to convert between domain entities and DTOs
func convertToOrderBookDTO(ob *entity.OrderBook) *dto.OrderBookDTO {
	return &dto.OrderBookDTO{
		ExchangeID: ob.ExchangeID(),
		Symbol:     ob.Symbol(),
		Bids:       convertToPriceLevelDTOs(ob.Bids()),
		Asks:       convertToPriceLevelDTOs(ob.Asks()),
		Timestamp:  ob.Timestamp(),
	}
}

func convertToOrderBookEntity(dto *dto.OrderBookDTO) *entity.OrderBook {
	ob := entity.NewOrderBook(dto.ExchangeID, dto.Symbol, dto.Timestamp)
	ob.UpdateBids(convertToPriceLevels(dto.Bids))
	ob.UpdateAsks(convertToPriceLevels(dto.Asks))
	return ob
}

func convertToPriceLevelDTOs(levels []entity.PriceLevel) []dto.PriceLevelDTO {
	dtos := make([]dto.PriceLevelDTO, len(levels))
	for i, level := range levels {
		dtos[i] = dto.PriceLevelDTO{
			Price:    level.Price.Value(),
			Quantity: level.Quantity.Value(),
		}
	}
	return dtos
}

func convertToPriceLevels(dtos []dto.PriceLevelDTO) []entity.PriceLevel {
	levels := make([]entity.PriceLevel, len(dtos))
	for i, dto := range dtos {
		price, _ := entity.NewPrice(dto.Price, "USDT")
		quantity, _ := entity.NewVolume(dto.Quantity, "")
		levels[i] = entity.PriceLevel{
			Price:    *price,
			Quantity: *quantity,
		}
	}
	return levels
}

func convertToTradeDTO(trade *entity.Trade) *dto.TradeDTO {
	return &dto.TradeDTO{
		ID:         trade.ID(),
		ExchangeID: trade.ExchangeID(),
		Symbol:     trade.Symbol(),
		Price:      trade.Price().Value(),
		Volume:     trade.Volume().Value(),
		TradeType:  string(trade.Type()),
		Timestamp:  trade.Timestamp(),
	}
}

// Ensure MarketDataService implements MarketDataUseCase interface
var _ input.MarketDataUseCase = (*MarketDataService)(nil)
