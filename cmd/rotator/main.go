package main

import (
	"banners-rotator/internal/config"
	"banners-rotator/internal/logger"
	"banners-rotator/internal/rotator"
	sqlstorage "banners-rotator/internal/storage/sql"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.dev.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.NewAppConfig(configFile)
	if err != nil {
		fmt.Printf("Critical app error: %v", err)
		os.Exit(1)
	}

	logg, err := logger.NewLogger(cfg.Logger.Level, []string{"stdout"})
	if err != nil {
		fmt.Printf("Critical app error: %v", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	logg.Info("getting storage...")
	s, err := getStorage(ctx, cfg)
	if err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1)
	}

	_ = rotator.NewApp(s)
}

func getStorage(ctx context.Context, cfg config.AppConfig) (rotator.Storage, error) {
	storage, err := sqlstorage.NewStorage(ctx, cfg.Storage.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("get storage: %w", err)
	}

	err = storage.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("get storage: %w", err)
	}

	return storage, nil
}
