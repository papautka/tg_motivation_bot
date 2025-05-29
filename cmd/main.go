package main

import (
	"log/slog"
	"os"
	"tg_motivation_bot/internal/adapters"
	"tg_motivation_bot/internal/config"
	"tg_motivation_bot/internal/sheduler"
	"tg_motivation_bot/internal/usecases"
)

func main() {
	app()
}

func app() {
	// 0. Создаем logger устанавливая его в slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 1. Подгружаем Config
	cfg := config.NewConfig()

	// 2. Подключаем adapter
	getQuoteAd := adapters.NewZenQuotesAdapter(cfg.QuoteAPIURL)
	translateQuoteAd := adapters.NewTranslateAdapter(cfg.TranslateAPIURL)
	sendMessageTgAd, err := adapters.NewTelegramAdapter(cfg.TelegramToken)

	if err != nil {
		logger.Info("Telegram adapter initialization failed")
		return
	}

	// 3. Usecases
	quoteFetcherUseCase := usecases.NewQuoteFetcher(getQuoteAd)
	translateFetcherUseCase := usecases.NewTranslateFetcher(translateQuoteAd)
	telegramFetcherUseCase := usecases.NewTelegramFetcher(sendMessageTgAd)

	// 4. Планировщик CRON
	cronExpr := "* * * * *" // Каждую минуту
	sheduler.InitScheduler(cronExpr, func() {
		// 4.1 получаем цитату на англосаксонском
		myQuote, errGetQuote := quoteFetcherUseCase.FetchFormattedQuote()
		if errGetQuote != nil {
			slog.Error("Ошибка при получении цитаты", slog.String("error", err.Error()))
			return
		}
		// 4.2 переводим на великий могучий
		mySlavicQuote, errTranslate := translateFetcherUseCase.FetchTranslated(myQuote, "en", "ru")
		if errTranslate != nil {
			slog.Error("Ошибка при переводе цитаты", slog.String("error", err.Error()))
			return
		}

		// 4.3 форматируем строку
		msg := usecases.FormatQuoteWithEmoji(mySlavicQuote.Text, mySlavicQuote.Author)

		// 4.4 отправляем в tg
		errSendMsg := telegramFetcherUseCase.FetchTelegram(cfg.TelegramChatId, msg)
		if errSendMsg != nil {
			slog.Error("Ошибка при отправке цитаты в telegram", slog.String("error", err.Error()))
			return
		} else {
			slog.Info("Цитата успешно отправлена")
		}
	})
	select {} // Блокируем main
}
