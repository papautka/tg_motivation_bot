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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	tgChatId, _ := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))
	return &Config{
		TelegramToken:   os.Getenv("TELEGRAM_TOKEN"),
		TelegramChatId:  tgChatId,
		QuoteAPIURL:     os.Getenv("QUOTE_API_URL"),
		TranslateAPIURL: os.Getenv("TRANSLATE_API_URL"),
	}
}
