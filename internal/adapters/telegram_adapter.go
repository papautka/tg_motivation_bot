package adapters

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type TelegramAdapter struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramAdapter(token string) (*TelegramAdapter, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to connect to telegram bot: %v", err)
		return nil, err
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &TelegramAdapter{bot: bot}, nil
}

// отправляет текстовое сообщение в указанный чат
func (t *TelegramAdapter) SendMessage(chatId int, message string) error {
	msg := tgbotapi.NewMessage(int64(chatId), message)
	msg.ParseMode = "Markdown"
	_, err := t.bot.Send(msg)
	return err
}
