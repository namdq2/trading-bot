package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"marketdata/internal/application/dto"
	"marketdata/internal/application/port/input"
)

type Consumer struct {
	reader            *kafka.Reader
	marketDataUseCase input.MarketDataUseCase
}

func NewConsumer(brokers []string, topic string, groupID string, useCase input.MarketDataUseCase) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{
		reader:            reader,
		marketDataUseCase: useCase,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				message, err := c.reader.ReadMessage(ctx)
				if err != nil {
					fmt.Printf("error reading message: %v\n", err)
					continue
				}

				if err := c.processMessage(ctx, message); err != nil {
					fmt.Printf("error processing message: %v\n", err)
				}
			}
		}
	}()

	return nil
}

func (c *Consumer) processMessage(ctx context.Context, message kafka.Message) error {
	var orderbook dto.OrderBookDTO
	if err := json.Unmarshal(message.Value, &orderbook); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	if err := c.marketDataUseCase.ProcessOrderBookUpdate(ctx, &orderbook); err != nil {
		return fmt.Errorf("failed to process orderbook update: %w", err)
	}

	return nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
