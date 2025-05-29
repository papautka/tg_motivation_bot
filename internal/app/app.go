package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tg_motivation_bot/internal/adapters"
	"tg_motivation_bot/internal/config"
	"tg_motivation_bot/internal/usecases"
	"time"
)

// Обновленная функция App с graceful shutdown
func App() {
	// 0. Создаем logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 1. Подгружаем Config
	cfg := config.NewConfig()
	if cfg == nil {
		logger.Error("Failed to load config")
		return
	}

	// 2. Подключаем adapters
	getQuoteAd := adapters.NewZenQuotesAdapter(cfg.QuoteAPIURL)
	translateQuoteAd := adapters.NewTranslateAdapter(cfg.TranslateAPIURL)
	sendMessageTgAd, err := adapters.NewTelegramAdapter(cfg.TelegramToken)

	if err != nil {
		logger.Error("Telegram adapter initialization failed", slog.String("error", err.Error()))
		return
	}

	// 3. Usecases
	quoteFetcherUseCase := usecases.NewQuoteFetcher(getQuoteAd)
	translateFetcherUseCase := usecases.NewTranslateFetcher(translateQuoteAd)
	telegramFetcherUseCase := usecases.NewTelegramFetcher(sendMessageTgAd)

	// 4. Запускаем бот
	go startTelegramBotLoop(sendMessageTgAd, telegramFetcherUseCase, quoteFetcherUseCase, translateFetcherUseCase)

	slog.Info("Bot started successfully")

	// 5. Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	slog.Info("Shutting down bot...")

	// Останавливаем бота
	sendMessageTgAd.Stop()

	// Даем время для завершения активных операций
	time.Sleep(2 * time.Second)

	slog.Info("Bot stopped gracefully")
}
