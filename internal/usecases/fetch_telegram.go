package usecases

import (
	"log/slog"
	"tg_motivation_bot/internal/interfaces"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramFetcher - use case для работы с Telegram
type TelegramFetcher struct {
	Provider interfaces.TelegramApi
}

// NewTelegramFetcher создает новый экземпляр TelegramFetcher
func NewTelegramFetcher(p interfaces.TelegramApi) *TelegramFetcher {
	return &TelegramFetcher{
		Provider: p,
	}
}

// StartBotLoop запускает основной цикл обработки обновлений
func (t *TelegramFetcher) StartBotLoopFetcher(handler func(string, int64)) {
	updates := t.Provider.StartBotLoop
	updates(handler)
	slog.Info("Telegram bot started and listening for updates...")
}

// FetchTelegram отправляет сообщение с reply клавиатурой (для обратной совместимости)
func (t *TelegramFetcher) FetchTelegram(chatId int64, msg string, keyboard *tgbotapi.ReplyKeyboardMarkup) error {
	var err error

	if keyboard != nil {
		err = t.Provider.SendMessageWithReplyKeyboard(chatId, msg, keyboard)
	} else {
		err = t.Provider.SendMessage(chatId, msg)
	}

	if err != nil {
		slog.Error("Error sending telegram message",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
			slog.String("message", truncateForLog(msg, 100)),
		)
	} else {
		slog.Info("Telegram message sent successfully",
			slog.Int64("chatId", chatId),
			slog.String("message", truncateForLog(msg, 50)),
		)
	}

	return err
}

// SendSimpleMessage отправляет простое текстовое сообщение
func (t *TelegramFetcher) SendSimpleMessage(chatId int64, msg string) error {
	err := t.Provider.SendMessage(chatId, msg)
	if err != nil {
		slog.Error("Error sending simple telegram message",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
		)
	}
	return err
}

// SendMessageWithInlineKeyboard отправляет сообщение с inline клавиатурой
func (t *TelegramFetcher) SendMessageWithInlineKeyboard(chatId int64, msg string, keyboard *tgbotapi.InlineKeyboardMarkup) error {
	err := t.Provider.SendMessageWithInlineKeyboard(chatId, msg, keyboard)
	if err != nil {
		slog.Error("Error sending telegram message with inline keyboard",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
		)
	}
	return err
}

// SendMessageWithDefaultKeyboard отправляет сообщение с предустановленной inline клавиатурой
func (t *TelegramFetcher) SendMessageWithDefaultKeyboard(chatId int64, msg string) error {
	err := t.Provider.SendMessageWithDefaultInlineKeyboard(chatId, msg)
	if err != nil {
		slog.Error("Error sending telegram message with default keyboard",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
		)
	}
	return err
}

// SendTypingIndicator отправляет индикатор "печатает..."
func (t *TelegramFetcher) SendTypingIndicator(chatId int64) error {
	err := t.Provider.SendTypingAction(chatId)
	if err != nil {
		slog.Warn("Error sending typing indicator",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
		)
	}
	return err
}

// SendPhoto отправляет фото с подписью
func (t *TelegramFetcher) SendPhoto(chatId int64, photoPath, caption string) error {
	err := t.Provider.SendPhoto(chatId, photoPath, caption)
	if err != nil {
		slog.Error("Error sending photo",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
			slog.String("photoPath", photoPath),
		)
	}
	return err
}

// SendDocument отправляет документ с подписью
func (t *TelegramFetcher) SendDocument(chatId int64, documentPath, caption string) error {
	err := t.Provider.SendDocument(chatId, documentPath, caption)
	if err != nil {
		slog.Error("Error sending document",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
			slog.String("documentPath", documentPath),
		)
	}
	return err
}

// EditMessage редактирует существующее сообщение
func (t *TelegramFetcher) EditMessage(chatId int64, messageId int, newText string) error {
	err := t.Provider.EditMessage(chatId, messageId, newText)
	if err != nil {
		slog.Error("Error editing message",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
			slog.Int("messageId", messageId),
		)
	}
	return err
}

// DeleteMessage удаляет сообщение
func (t *TelegramFetcher) DeleteMessage(chatId int64, messageId int) error {
	err := t.Provider.DeleteMessage(chatId, messageId)
	if err != nil {
		slog.Error("Error deleting message",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
			slog.Int("messageId", messageId),
		)
	}
	return err
}

// GetBotInfo получает информацию о боте
func (t *TelegramFetcher) GetBotInfo() (*tgbotapi.User, error) {
	botInfo, err := t.Provider.GetMe()
	if err != nil {
		slog.Error("Error getting bot info", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("Bot info retrieved",
		slog.String("username", botInfo.UserName),
		slog.String("firstName", botInfo.FirstName),
	)

	return botInfo, nil
}

// GetChatMember получает информацию о участнике чата
func (t *TelegramFetcher) GetChatMember(chatId, userId int64) (*tgbotapi.ChatMember, error) {
	member, err := t.Provider.GetChatMember(chatId, userId)
	if err != nil {
		slog.Error("Error getting chat member",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
			slog.Int64("userId", userId),
		)
		return nil, err
	}
	return member, nil
}

// SetBotCommands устанавливает команды бота
func (t *TelegramFetcher) SetBotCommands(commands []tgbotapi.BotCommand) error {
	// 1. Устанавливаем команды для бота
	err := t.Provider.SetCommands(commands)
	if err != nil {
		slog.Error("Error setting bot commands",
			slog.String("error", err.Error()),
			slog.Int("commandsCount", len(commands)),
		)
		return err
	}

	slog.Info("Bot commands set successfully", slog.Int("count", len(commands)))
	return nil
}

// IsPrivateChat проверяет, является ли чат приватным
func (t *TelegramFetcher) IsPrivateChat(chatId int64) bool {
	return t.Provider.IsPrivateChat(chatId)
}

// IsGroupChat проверяет, является ли чат групповым
func (t *TelegramFetcher) IsGroupChat(chatId int64) bool {
	return t.Provider.IsGroupChat(chatId)
}

// StopBot останавливает бота
func (t *TelegramFetcher) StopBot() {
	t.Provider.Stop()
	slog.Info("Bot stopped via TelegramFetcher")
}

// NotifyError отправляет уведомление об ошибке с кнопками
func (t *TelegramFetcher) NotifyError(chatId int64, errorMsg string) error {
	fullMsg := "❌ " + errorMsg + "\n\nПопробуйте еще раз:"
	return t.SendMessageWithDefaultKeyboard(chatId, fullMsg)
}

// NotifySuccess отправляет уведомление об успехе с кнопками
func (t *TelegramFetcher) NotifySuccess(chatId int64, successMsg string) error {
	fullMsg := "✅ " + successMsg + "\n\nВыберите действие:"
	return t.SendMessageWithDefaultKeyboard(chatId, fullMsg)
}

// SendQuoteWithKeyboard отправляет отформатированную цитату с клавиатурой
func (t *TelegramFetcher) SendQuoteWithKeyboard(chatId int64, text, author string) error {
	formattedQuote := FormatQuoteMessage(text, author)
	return t.SendMessageWithDefaultKeyboard(chatId, formattedQuote)
}

// FormatQuoteMessage форматирует цитату для отправки
func FormatQuoteMessage(text, author string) string {
	return "💬 \"" + text + "\"\n\n— " + author
}

// truncateForLog обрезает строку для логирования
func truncateForLog(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Методы для работы с различными типами уведомлений

// SendWelcomeMessage отправляет приветственное сообщение
func (t *TelegramFetcher) SendWelcomeMessage(chatId int64) error {
	welcomeText := `🎯 Добро пожаловать в Quote Bot!

Этот бот поможет вам получить вдохновляющие цитаты на русском или английском языке.

Выберите язык для получения цитаты:`

	return t.SendMessageWithDefaultKeyboard(chatId, welcomeText)
}

// SendHelpMessage отправляет справочное сообщение
func (t *TelegramFetcher) SendHelpMessage(chatId int64) error {
	helpText := `ℹ️ Помощь по использованию бота:

🔹 /start - начать работу с ботом
🔹 /quote - получить цитату
🔹 /help - показать эту справку

Просто нажимайте на кнопки для выбора языка цитаты!`

	return t.SendMessageWithDefaultKeyboard(chatId, helpText)
}
