package adapters

import (
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramAdapter - адаптер для работы с Telegram Bot API
type TelegramAdapter struct {
	bot *tgbotapi.BotAPI
}

// Предопределенные клавиатуры
var (
	// Inline клавиатура с выбором языка
	InlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "ru"),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "en"),
		),
	)

	// Reply клавиатура (если понадобится)
	DefaultKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🇷🇺 Русский"),
			tgbotapi.NewKeyboardButton("🇬🇧 English"),
		),
	)
)

// NewTelegramAdapter создает новый экземпляр TelegramAdapter
func NewTelegramAdapter(token string) (*TelegramAdapter, error) {
	if token == "" {
		return nil, fmt.Errorf("telegram token cannot be empty")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	// Включаем debug режим для разработки (можно убрать в продакшене)
	bot.Debug = false

	slog.Info("Telegram bot initialized", slog.String("username", bot.Self.UserName))

	return &TelegramAdapter{
		bot: bot,
	}, nil
}

// StartBotLoop запускает основной цикл обработки обновлений
func (t *TelegramAdapter) StartBotLoop(handler func(string, int64)) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	slog.Info("Telegram bot started and listening for updates...")

	for update := range updates {
		// Обрабатываем текстовые сообщения
		if update.Message != nil {
			go t.handleMessage(update.Message, handler)
		}

		// Обрабатываем нажатия на inline кнопки
		if update.CallbackQuery != nil {
			go t.handleCallbackQuery(update.CallbackQuery, handler)
		}
	}
}

// handleMessage обрабатывает текстовые сообщения
func (t *TelegramAdapter) handleMessage(message *tgbotapi.Message, handler func(string, int64)) {
	chatID := message.Chat.ID
	text := message.Text

	slog.Info("Received message",
		slog.String("text", text),
		slog.Int64("chatId", chatID),
		slog.String("username", message.From.UserName),
	)

	// Вызываем обработчик
	handler(text, chatID)
}

// handleCallbackQuery обрабатывает нажатия на inline кнопки
func (t *TelegramAdapter) handleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, handler func(string, int64)) {
	chatID := callbackQuery.Message.Chat.ID
	data := callbackQuery.Data

	slog.Info("Received callback query",
		slog.String("data", data),
		slog.Int64("chatId", chatID),
		slog.String("username", callbackQuery.From.UserName),
	)

	// Отвечаем на callback query (убирает индикатор загрузки)
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := t.bot.Request(callback); err != nil {
		slog.Error("Failed to answer callback query", slog.String("error", err.Error()))
	}

	// Вызываем обработчик
	handler(data, chatID)
}

// SendMessage отправляет простое текстовое сообщение
func (t *TelegramAdapter) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	return t.sendMessage(msg)
}

// SendMessageWithReplyKeyboard отправляет сообщение с reply клавиатурой
func (t *TelegramAdapter) SendMessageWithReplyKeyboard(chatID int64, text string, keyboard *tgbotapi.ReplyKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	if keyboard != nil {
		msg.ReplyMarkup = *keyboard
	}
	return t.sendMessage(msg)
}

// SendMessageWithInlineKeyboard отправляет сообщение с inline клавиатурой
func (t *TelegramAdapter) SendMessageWithInlineKeyboard(chatID int64, text string, keyboard *tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	if keyboard != nil {
		msg.ReplyMarkup = *keyboard
	}
	return t.sendMessage(msg)
}

// SendMessageWithDefaultInlineKeyboard отправляет сообщение с предустановленной inline клавиатурой
func (t *TelegramAdapter) SendMessageWithDefaultInlineKeyboard(chatID int64, text string) error {
	return t.SendMessageWithInlineKeyboard(chatID, text, &InlineKeyboard)
}

// SendPhoto отправляет фото
func (t *TelegramAdapter) SendPhoto(chatID int64, photoPath string, caption string) error {
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(photoPath))
	photo.Caption = caption

	_, err := t.bot.Send(photo)
	if err != nil {
		slog.Error("Failed to send photo",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatID),
		)
		return fmt.Errorf("failed to send photo: %w", err)
	}

	slog.Info("Photo sent successfully", slog.Int64("chatId", chatID))
	return nil
}

// SendDocument отправляет документ
func (t *TelegramAdapter) SendDocument(chatID int64, documentPath string, caption string) error {
	document := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(documentPath))
	document.Caption = caption

	_, err := t.bot.Send(document)
	if err != nil {
		slog.Error("Failed to send document",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatID),
		)
		return fmt.Errorf("failed to send document: %w", err)
	}

	slog.Info("Document sent successfully", slog.Int64("chatId", chatID))
	return nil
}

// SendTypingAction отправляет индикатор "печатает..."
func (t *TelegramAdapter) SendTypingAction(chatID int64) error {
	action := tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping)
	_, err := t.bot.Request(action)
	if err != nil {
		slog.Error("Failed to send typing action",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatID),
		)
		return fmt.Errorf("failed to send typing action: %w", err)
	}
	return nil
}

// EditMessage редактирует существующее сообщение
func (t *TelegramAdapter) EditMessage(chatID int64, messageID int, newText string) error {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, newText)
	_, err := t.bot.Send(edit)
	if err != nil {
		slog.Error("Failed to edit message",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatID),
			slog.Int("messageId", messageID),
		)
		return fmt.Errorf("failed to edit message: %w", err)
	}

	slog.Info("Message edited successfully",
		slog.Int64("chatId", chatID),
		slog.Int("messageId", messageID),
	)
	return nil
}

// DeleteMessage удаляет сообщение
func (t *TelegramAdapter) DeleteMessage(chatID int64, messageID int) error {
	delete := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := t.bot.Request(delete)
	if err != nil {
		slog.Error("Failed to delete message",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatID),
			slog.Int("messageId", messageID),
		)
		return fmt.Errorf("failed to delete message: %w", err)
	}

	slog.Info("Message deleted successfully",
		slog.Int64("chatId", chatID),
		slog.Int("messageId", messageID),
	)
	return nil
}

// GetMe возвращает информацию о боте
func (t *TelegramAdapter) GetMe() (*tgbotapi.User, error) {
	return &t.bot.Self, nil
}

// GetChatMember получает информацию о участнике чата
func (t *TelegramAdapter) GetChatMember(chatID int64, userID int64) (*tgbotapi.ChatMember, error) {
	member, err := t.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get chat member: %w", err)
	}
	return &member, nil
}

// SetCommands устанавливает команды бота
func (t *TelegramAdapter) SetCommands(commands []tgbotapi.BotCommand) error {
	config := tgbotapi.NewSetMyCommands(commands...)
	_, err := t.bot.Request(config)
	if err != nil {
		slog.Error("Failed to set bot commands", slog.String("error", err.Error()))
		return fmt.Errorf("failed to set bot commands: %w", err)
	}

	slog.Info("Bot commands set successfully", slog.Int("count", len(commands)))
	return nil
}

// sendMessage - внутренний метод для отправки сообщений с retry логикой
func (t *TelegramAdapter) sendMessage(msg tgbotapi.MessageConfig) error {
	const maxRetries = 3
	const retryDelay = time.Second

	var err error
	for i := 0; i < maxRetries; i++ {
		_, err = t.bot.Send(msg)
		if err == nil {
			slog.Info("Message sent successfully",
				slog.Int64("chatId", msg.ChatID),
				slog.String("text", truncateString(msg.Text, 50)),
			)
			return nil
		}

		slog.Warn("Failed to send message, retrying...",
			slog.String("error", err.Error()),
			slog.Int("attempt", i+1),
			slog.Int64("chatId", msg.ChatID),
		)

		if i < maxRetries-1 {
			time.Sleep(retryDelay * time.Duration(i+1))
		}
	}

	slog.Error("Failed to send message after all retries",
		slog.String("error", err.Error()),
		slog.Int64("chatId", msg.ChatID),
	)
	return fmt.Errorf("failed to send message after %d attempts: %w", maxRetries, err)
}

// truncateString обрезает строку до указанной длины для логирования
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Дополнительные методы для удобства

// IsPrivateChat проверяет, является ли чат приватным
func (t *TelegramAdapter) IsPrivateChat(chatID int64) bool {
	return chatID > 0
}

// IsGroupChat проверяет, является ли чат групповым
func (t *TelegramAdapter) IsGroupChat(chatID int64) bool {
	return chatID < 0
}

// Stop останавливает бота (полезно для graceful shutdown)
func (t *TelegramAdapter) Stop() {
	t.bot.StopReceivingUpdates()
	slog.Info("Telegram bot stopped")
}
