package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"marketdata/internal/application/dto"
)

type Publisher struct {
	writer *kafka.Writer
	topic  string
}

func NewPublisher(brokers []string, topic string) *Publisher {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	return &Publisher{
		writer: writer,
		topic:  topic,
	}
}

func (p *Publisher) PublishOrderBookUpdate(ctx context.Context, orderbook *dto.OrderBookDTO) error {
	data, err := json.Marshal(orderbook)
	if err != nil {
		return fmt.Errorf("failed to marshal orderbook: %w", err)
	}

	key := fmt.Sprintf("%s-%s", orderbook.ExchangeID, orderbook.Symbol)
	message := kafka.Message{
		Key:   []byte(key),
		Value: data,
	}

	if err := p.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
