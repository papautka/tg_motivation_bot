package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	TelegramToken   string
	TelegramChatId  int64
	QuoteAPIURL     string
	TranslateAPIURL string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	chatIdStr := os.Getenv("TELEGRAM_CHAT_ID")
	log.Printf("Telegram token: %s", token)       // для дебага
	log.Printf("Telegram chat ID: %s", chatIdStr) // для дебага

	tgChatId, _ := strconv.Atoi(chatIdStr)
	return &Config{
		TelegramToken:   token,
		TelegramChatId:  int64(tgChatId),
		QuoteAPIURL:     os.Getenv("QUOTE_API_URL"),
		TranslateAPIURL: os.Getenv("TRANSLATE_API_URL"),
	}
}
