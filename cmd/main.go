package main

import (
	"fmt"
	"tg_motivation_bot/internal/config"
)

func main() {
	app()
}

func app() {
	// 1. Подгружаем Config
	cfg := config.NewConfig()
	fmt.Printf("%+v\n", cfg)
}
