package timescale

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"marketdata/internal/domain/entity"
	"marketdata/internal/domain/valueobject"
)

type TradeRepository struct {
	db *sql.DB
}

func NewTradeRepository(db *sql.DB) *TradeRepository {
	return &TradeRepository{
		db: db,
	}
}

func (r *TradeRepository) StoreTrade(ctx context.Context, trade *entity.Trade) error {
	query := `
		INSERT INTO trades (
			id, exchange_id, symbol, price, volume, trade_type, timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx,
		query,
		trade.ID(),
		trade.ExchangeID(),
		trade.Symbol(),
		trade.Price().Value(),
		trade.Volume().Value(),
		trade.Type(),
		trade.Timestamp(),
	)

	if err != nil {
		return fmt.Errorf("failed to store trade: %w", err)
	}

	return nil
}

func (r *TradeRepository) GetTradesBySymbol(ctx context.Context, symbol string, limit int) ([]*entity.Trade, error) {
	query := `
		SELECT id, exchange_id, symbol, price, volume, trade_type, timestamp
		FROM trades
		WHERE symbol = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, symbol, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query trades: %w", err)
	}
	defer rows.Close()

	var trades []*entity.Trade
	for rows.Next() {
		var (
			id         string
			exchangeID string
			symbol     string
			price      float64
			volume     float64
			tradeType  string
			timestamp  time.Time
		)

		if err := rows.Scan(&id, &exchangeID, &symbol, &price, &volume, &tradeType, &timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan trade row: %w", err)
		}

		priceVO, err := valueobject.NewPrice(price, "USDT")
		if err != nil {
			return nil, fmt.Errorf("failed to create price value object: %w", err)
		}

		volumeVO, err := valueobject.NewVolume(volume, symbol)
		if err != nil {
			return nil, fmt.Errorf("failed to create volume value object: %w", err)
		}

		trade := entity.NewTrade(
			id,
			exchangeID,
			symbol,
			*priceVO,
			*volumeVO,
			entity.TradeType(tradeType),
			timestamp,
		)

		trades = append(trades, trade)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trade rows: %w", err)
	}

	return trades, nil
}
