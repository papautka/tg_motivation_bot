package interfaces

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Quote struct {
	Text   string
	Author string
}

type QuoteProvider interface {
	GetQuote() (*Quote, error)
}

type Translator interface {
	Translate(quote *Quote, fromLang, toLang string) (*Quote, error)
}

// TelegramApi - интерфейс для работы с Telegram API
type TelegramApi interface {
	// Основные методы отправки сообщений
	SendMessage(chatID int64, text string) error
	SendMessageWithReplyKeyboard(chatID int64, text string, keyboard *tgbotapi.ReplyKeyboardMarkup) error
	SendMessageWithInlineKeyboard(chatID int64, text string, keyboard *tgbotapi.InlineKeyboardMarkup) error
	SendMessageWithDefaultInlineKeyboard(chatID int64, text string) error

	// Дополнительные методы
	SendTypingAction(chatID int64) error
	SendPhoto(chatID int64, photoPath string, caption string) error
	SendDocument(chatID int64, documentPath string, caption string) error

	// Управление сообщениями
	EditMessage(chatID int64, messageID int, newText string) error
	DeleteMessage(chatID int64, messageID int) error

	// Информационные методы
	GetMe() (*tgbotapi.User, error)
	GetChatMember(chatID int64, userID int64) (*tgbotapi.ChatMember, error)

	// Управление ботом
	SetCommands(commands []tgbotapi.BotCommand) error
	Stop()

	// Утилиты
	IsPrivateChat(chatID int64) bool
	IsGroupChat(chatID int64) bool
}
