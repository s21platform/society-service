package main

import (
	"fmt"
	"net"
	"os"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/society-service/internal/infra"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	society "github.com/s21platform/society-proto/society-proto"

	"github.com/s21platform/society-service/internal/config"
	db "github.com/s21platform/society-service/internal/repository/postgres"
	"github.com/s21platform/society-service/internal/service"
)

func main() {
	// чтение конфига
	cfg := config.MustLoad()

	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	dbRepo, err := db.New(cfg)

	if err != nil {
		logger.Error(fmt.Sprintf("failed to db.New: %v", err))
		os.Exit(1)
	}
	defer dbRepo.Close()

	server := service.New(dbRepo)
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.Verifcation,
		),
		grpc.ChainUnaryInterceptor(infra.Logger(logger)),
	)
	society.RegisterSocietyServiceServer(s, server)

	logger.Info(fmt.Sprintf("starting server %v", cfg.Service.Port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to cannnot listen port: %s; Error: %s", cfg.Service.Port, err))
	}
	if err := s.Serve(lis); err != nil {
		logger.Error(fmt.Sprintf("failed to cannnot start service: %s; Error: %s", cfg.Service.Port, err))
	}
}
