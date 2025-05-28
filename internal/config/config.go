package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	TelegramToken   string
	TelegramChatId  string
	QuoteAPIURL     string
	TranslateAPIURL string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		TelegramToken:   os.Getenv("TELEGRAM_TOKEN"),
		TelegramChatId:  os.Getenv("TELEGRAM_CHAT_ID"),
		QuoteAPIURL:     os.Getenv("QUOTE_API_URL"),
		TranslateAPIURL: os.Getenv("TRANSLATE_API_URL"),
	}
}
