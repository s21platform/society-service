package main

import (
	"github.com/s21platform/society-service/internal/config"
	db "github.com/s21platform/society-service/internal/repository/postgres"
	"log"
	"os"
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
}
