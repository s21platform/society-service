package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	society "github.com/s21platform/society-proto/society-proto"

	"github.com/s21platform/society-service/internal/config"
	db "github.com/s21platform/society-service/internal/repository/postgres"
	"github.com/s21platform/society-service/internal/rpc"
)

func main() {
	// чтение конфига
	cfg := config.MustLoad()

	dbRepo, err := db.New(cfg)

	if err != nil {
		log.Printf("db.New: %v", err)
		os.Exit(1)
	}
	defer dbRepo.Close()

	server := rpc.New(dbRepo)
	s := grpc.NewServer()
	society.RegisterSocietyServiceServer(s, server)

	log.Println("starting server", cfg.Service.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("Cannnot listen port: %s; Error: %s", cfg.Service.Port, err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannnot start rpc: %s; Error: %s", cfg.Service.Port, err)
	}
}
