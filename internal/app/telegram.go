package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"tg_motivation_bot/internal/adapters"
	"tg_motivation_bot/internal/usecases"
)

// startTelegramBotLoop запускает цикл обработки сообщений
func startTelegramBotLoop(
	bot *adapters.TelegramAdapter,
	tgFetcher *usecases.TelegramFetcher,
	quoteFetcher *usecases.QuoteFetcher,
	translateFetcher *usecases.TranslateFetcher,
) {
	// Устанавливаем команды бота
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Начать работу с ботом"},
		{Command: "help", Description: "Помощь"},
		{Command: "quote", Description: "Получить цитату"},
	}

	if err := bot.SetCommands(commands); err != nil {
		slog.Error("Failed to set bot commands", slog.String("error", err.Error()))
	}

	// Запускаем обработчик сообщений
	bot.StartBotLoop(func(command string, chatId int64) {
		handleBotCommand(command, chatId, bot, quoteFetcher, translateFetcher)
	})
}

// handleBotCommand обрабатывает команды бота
func handleBotCommand(
	command string,
	chatId int64,
	bot *adapters.TelegramAdapter,
	quoteFetcher *usecases.QuoteFetcher,
	translateFetcher *usecases.TranslateFetcher,
) {
	slog.Info("Processing command",
		slog.String("command", command),
		slog.Int64("chatId", chatId),
	)

	// Отправляем индикатор "печатает..."
	if err := bot.SendTypingAction(chatId); err != nil {
		slog.Warn("Failed to send typing action", slog.String("error", err.Error()))
	}

	switch command {
	case "/start":
		handleStartCommand(chatId, bot)
	case "/help":
		handleHelpCommand(chatId, bot)
	case "/quote":
		// Показываем меню выбора языка
		bot.SendMessageWithDefaultInlineKeyboard(chatId, "Выберите язык для цитаты:")
	case "ru":
		handleRussianQuote(chatId, bot, quoteFetcher, translateFetcher)
	case "en":
		handleEnglishQuote(chatId, bot, quoteFetcher)
	default:
		handleUnknownCommand(chatId, bot, command)
	}
}

// handleStartCommand обрабатывает команду /start
func handleStartCommand(chatId int64, bot *adapters.TelegramAdapter) {
	welcomeText := `🎯 Добро пожаловать в Quote Bot!

Этот бот поможет вам получить вдохновляющие цитаты на русском или английском языке.

Выберите язык для получения цитаты:`

	bot.SendMessageWithDefaultInlineKeyboard(chatId, welcomeText)
}

// handleHelpCommand обрабатывает команду /help
func handleHelpCommand(chatId int64, bot *adapters.TelegramAdapter) {
	helpText := `ℹ️ Помощь по использованию бота:

🔹 /start - начать работу с ботом
🔹 /quote - получить цитату
🔹 /help - показать эту справку

Просто нажимайте на кнопки для выбора языка цитаты!`

	bot.SendMessageWithDefaultInlineKeyboard(chatId, helpText)
}

// handleRussianQuote обрабатывает запрос цитаты на русском языке
func handleRussianQuote(
	chatId int64,
	bot *adapters.TelegramAdapter,
	quoteFetcher *usecases.QuoteFetcher,
	translateFetcher *usecases.TranslateFetcher,
) {
	// Получаем цитату
	quote, err := quoteFetcher.FetchFormattedQuote()
	if err != nil {
		slog.Error("Failed to fetch quote",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
		)
		bot.SendMessageWithDefaultInlineKeyboard(chatId, "❌ Ошибка при получении цитаты. Попробуйте еще раз.")
		return
	}

	// Переводим на русский
	translated, err := translateFetcher.FetchTranslated(quote, "en", "ru")
	if err != nil {
		slog.Error("Failed to translate quote",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
		)
		bot.SendMessageWithDefaultInlineKeyboard(chatId, "❌ Ошибка при переводе цитаты. Попробуйте еще раз.")
		return
	}

	// Форматируем и отправляем
	formattedQuote := usecases.FormatQuoteWithEmoji(translated.Text, translated.Author)
	bot.SendMessageWithDefaultInlineKeyboard(chatId, formattedQuote)

	slog.Info("Russian quote sent successfully", slog.Int64("chatId", chatId))
}

// handleEnglishQuote обрабатывает запрос цитаты на английском языке
func handleEnglishQuote(
	chatId int64,
	bot *adapters.TelegramAdapter,
	quoteFetcher *usecases.QuoteFetcher,
) {
	// Получаем цитату
	quote, err := quoteFetcher.FetchFormattedQuote()
	if err != nil {
		slog.Error("Failed to fetch quote",
			slog.String("error", err.Error()),
			slog.Int64("chatId", chatId),
		)
		bot.SendMessageWithDefaultInlineKeyboard(chatId, "❌ Error fetching quote. Please try again.")
		return
	}

	// Форматируем и отправляем
	formattedQuote := usecases.FormatQuoteWithEmoji(quote.Text, quote.Author)
	bot.SendMessageWithDefaultInlineKeyboard(chatId, formattedQuote)

	slog.Info("English quote sent successfully", slog.Int64("chatId", chatId))
}

// handleUnknownCommand обрабатывает неизвестные команды
func handleUnknownCommand(chatId int64, bot *adapters.TelegramAdapter, command string) {
	unknownText := `❓ Неизвестная команда: ` + command + `Используйте /help для получения списка доступных команд. Или выберите действие:`

	bot.SendMessageWithDefaultInlineKeyboard(chatId, unknownText)

	slog.Info("Unknown command received",
		slog.String("command", command),
		slog.Int64("chatId", chatId),
	)
}
