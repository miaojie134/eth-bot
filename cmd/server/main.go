package main

import (
	"log"
	"net/http"

	"github.com/qqqq/eth-trading-system/internal/api"
	"github.com/qqqq/eth-trading-system/internal/config"
	"github.com/qqqq/eth-trading-system/internal/services"
	"github.com/qqqq/eth-trading-system/internal/storage"
	"github.com/qqqq/eth-trading-system/internal/utils"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// 初始化日志系统
	err = utils.InitLogger(cfg.LogDir, cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// 使用新的日志系统
	utils.Log.Info("Starting ETH Trading System")

	db, err := storage.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	alpacaService := services.NewAlpacaService(cfg.AlpacaAPIKey, cfg.AlpacaAPISecret)

	// 创建并启动数据收集服务
	dataCollectionService := services.NewDataCollectionService(alpacaService, db)
	dataCollectionService.Start()

	handler := api.NewHandler(alpacaService, dataCollectionService)

	http.HandleFunc("/", handler.IndexHandler)
	http.HandleFunc("/api/price", handler.GetLatestPrice)
	http.HandleFunc("/api/historical", handler.GetHistoricalData)

	utils.Log.Infof("Server starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, nil); err != nil {
		utils.Log.Fatalf("Failed to start server: %v", err)
	}
}
