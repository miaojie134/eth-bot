package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// 创建必要的表
	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	// 创建market_data表，添加timeframe字段，并创建唯一索引
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS market_data (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            symbol TEXT NOT NULL,
            timeframe TEXT NOT NULL,
            open REAL NOT NULL,
            high REAL NOT NULL,
            low REAL NOT NULL,
            close REAL NOT NULL,
            volume REAL NOT NULL,
            timestamp DATETIME NOT NULL,
            trade_count INTEGER NOT NULL,
            vwap REAL NOT NULL,
            UNIQUE(symbol, timeframe, timestamp)
        )
    `)
	if err != nil {
		return err
	}

	// 创建latest_price表，并创建唯一索引
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS latest_price (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            symbol TEXT NOT NULL,
            open REAL NOT NULL,
            high REAL NOT NULL,
            low REAL NOT NULL,
            close REAL NOT NULL,
            volume REAL NOT NULL,
            timestamp DATETIME NOT NULL,
            trade_count INTEGER NOT NULL,
            vwap REAL NOT NULL,
            UNIQUE(symbol)
        )
    `)
	if err != nil {
		return err
	}

	return nil
}
