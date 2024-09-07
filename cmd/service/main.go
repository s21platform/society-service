package main

import (
	"fmt"
	"github.com/s21platform/society-service/internal/config"
)

func main() {
	// чтение конфига
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
