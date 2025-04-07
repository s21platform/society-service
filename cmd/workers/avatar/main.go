package main

import (
	"context"
	kafkalib "github.com/s21platform/kafka-lib"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"
	"github.com/s21platform/society-service/internal/config"
	"github.com/s21platform/society-service/internal/databus/avatar"
	db "github.com/s21platform/society-service/internal/repository/postgres"
	"os"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)
	dbRepo, err := db.New(cfg)

	if err != nil {
		logger.Error("failed to postgres.New")
		os.Exit(1)
	}
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "society", cfg.Platform.Env)
	if err != nil {
		logger.Error("failed to connect grafit")
		os.Exit(1)
	}

	ctx := context.WithValue(context.Background(), config.KeyMetrics, metrics)

	//Consumer
	newSocietyConsumer, err := kafkalib.NewConsumer(cfg.Kafka.Server, cfg.Kafka.SocietyNewAvatar, metrics)
	if err != nil {
		logger.Error("failed to create kafka consumer")
		os.Exit(1)
	}

	NewSocietyHandler := avatar.New(dbRepo)

	newSocietyConsumer.RegisterHandler(ctx, NewSocietyHandler.Handler)

	<-ctx.Done()
}
