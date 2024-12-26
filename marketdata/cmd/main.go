package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"marketdata/config"
	"marketdata/internal/application/service"
	"marketdata/internal/exchange"
	"marketdata/internal/repository"
	"marketdata/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()
	defer log.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize market data service
	svc, err := service.NewMarketDataService(
		ctx,
		cfg,
		log,
		repository.NewOrderBookRepository(redisClient),
		repository.NewTradeRepository(dbClient),
		exchange.NewExchangeManager(cfg.Exchange),
	)
	if err != nil {
		log.Fatal("failed to create market data service", err)
	}

	// Start the service
	if err := svc.Start(ctx); err != nil {
		log.Fatal("failed to start market data service", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("shutting down market data service...")
	if err := svc.Stop(); err != nil {
		log.Error("error during shutdown", err)
	}
}
