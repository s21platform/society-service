package main

import (
	"fmt"
	"log"
	"net"
	"os"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/society-service/internal/infra"

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
		log.Printf("db.New: %v", err)
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

	log.Println("starting server", cfg.Service.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("Cannnot listen port: %s; Error: %s", cfg.Service.Port, err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannnot start service: %s; Error: %s", cfg.Service.Port, err)
	}
}
