package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	TelegramToken   string
	TelegramChatId  int
	QuoteAPIURL     string
	TranslateAPIURL string
}

func NewConfig() *Config {
	// Попытаться загрузить .env файл, но не ругаться, если его нет
	// (Например, в продакшене .env не будет — и это нормально)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	tgChatId, _ := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))
	return &Config{
		TelegramToken:   os.Getenv("TELEGRAM_TOKEN"),
		TelegramChatId:  tgChatId,
		QuoteAPIURL:     os.Getenv("QUOTE_API_URL"),
		TranslateAPIURL: os.Getenv("TRANSLATE_API_URL"),
	}
}
