package storage

import (
	"database/sql"
	"time"

	"github.com/qqqq/eth-trading-system/internal/models"
)

type DataRepository interface {
	StoreLatestPrice(bar *models.Bar) error
	StoreHistoricalData(bars []models.Bar, timeframe string) error
	GetLatestPrice() (*models.Bar, error)
	GetHistoricalData(timeframe string, start, end time.Time) ([]models.Bar, error)
}

type SQLiteDataRepository struct {
	db *sql.DB
}

func NewDataRepository(db *sql.DB) DataRepository {
	return &SQLiteDataRepository{db: db}
}

func (r *SQLiteDataRepository) StoreLatestPrice(bar *models.Bar) error {
	_, err := r.db.Exec(`
		INSERT INTO latest_price (symbol, open, high, low, close, volume, timestamp, trade_count, vwap)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(symbol) DO UPDATE SET
		open=excluded.open, high=excluded.high, low=excluded.low, close=excluded.close, volume=excluded.volume, timestamp=excluded.timestamp, trade_count=excluded.trade_count, vwap=excluded.vwap
	`, "ETH/USD", bar.Open, bar.High, bar.Low, bar.Close, bar.Volume, bar.Timestamp, bar.TradeCount, bar.VWAP)
	return err
}

func (r *SQLiteDataRepository) StoreHistoricalData(bars []models.Bar, timeframe string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO market_data (symbol, timeframe, open, high, low, close, volume, timestamp, trade_count, vwap)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(symbol, timeframe, timestamp) DO UPDATE SET
		open=excluded.open, high=excluded.high, low=excluded.low, close=excluded.close, volume=excluded.volume, trade_count=excluded.trade_count, vwap=excluded.vwap
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, bar := range bars {
		_, err := stmt.Exec("ETH/USD", timeframe, bar.Open, bar.High, bar.Low, bar.Close, bar.Volume, bar.Timestamp, bar.TradeCount, bar.VWAP)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteDataRepository) GetLatestPrice() (*models.Bar, error) {
	var bar models.Bar
	err := r.db.QueryRow(`
		SELECT open, high, low, close, volume, timestamp, trade_count, vwap
		FROM latest_price
		WHERE symbol = 'ETH/USD'
	`).Scan(&bar.Open, &bar.High, &bar.Low, &bar.Close, &bar.Volume, &bar.Timestamp, &bar.TradeCount, &bar.VWAP)

	if err != nil {
		return nil, err
	}

	return &bar, nil
}

func (r *SQLiteDataRepository) GetHistoricalData(timeframe string, start, end time.Time) ([]models.Bar, error) {
	rows, err := r.db.Query(`
		SELECT open, high, low, close, volume, timestamp, trade_count, vwap
		FROM market_data
		WHERE symbol = 'ETH/USD' AND timeframe = ? AND timestamp BETWEEN ? AND ?
		ORDER BY timestamp ASC
	`, timeframe, start, end)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bars []models.Bar
	for rows.Next() {
		var bar models.Bar
		err := rows.Scan(&bar.Open, &bar.High, &bar.Low, &bar.Close, &bar.Volume, &bar.Timestamp, &bar.TradeCount, &bar.VWAP)
		if err != nil {
			return nil, err
		}
		bars = append(bars, bar)
	}

	return bars, nil
}
