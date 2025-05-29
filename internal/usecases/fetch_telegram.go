package usecases

import (
	"log"
	"tg_motivation_bot/internal/interfaces"
)

type TelegramFetcher struct {
	Provider interfaces.TelegramApi
}

func NewTelegramFetcher(p interfaces.TelegramApi) *TelegramFetcher {
	return &TelegramFetcher{
		Provider: p,
	}
}

func (t *TelegramFetcher) FetchTelegram(chatId int, msg string) error {
	err := t.Provider.SendMessage(chatId, msg)
	if err != nil {
		log.Printf("(FetchTelegram) Error sending message: %v", err)
	}
	return err
}
