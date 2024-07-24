package services

import (
	"time"

	"github.com/qqqq/eth-trading-system/internal/models"
)

type AlpacaServiceInterface interface {
	GetHistoricalBars(symbol, timeframe, start, end string, limit int, pageToken string) ([]models.AlpacaBar, string, error)
	GetLatestBar() (*models.AlpacaBar, error)
}

type DataCollectionServiceInterface interface {
	Start()
	GetLatestPrice() (*models.Bar, error)
	GetHistoricalData(timeframe string, start, end time.Time) ([]models.Bar, error)
}

type StorageInterface interface {
	StoreHistoricalData(bar models.Bar, timeframe string) error
	StoreLatestPrice(bar models.Bar) error
	GetLatestPrice() (*models.Bar, error)
	GetHistoricalData(timeframe string, start, end time.Time) ([]models.Bar, error)
}
