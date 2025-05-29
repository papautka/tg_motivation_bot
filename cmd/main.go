package main

import (
	"fmt"
	"tg_motivation_bot/internal/adapters"
	"tg_motivation_bot/internal/config"
	"tg_motivation_bot/internal/usecases"
)

func main() {
	app()
}

func app() {
	// 1. Подгружаем Config
	cfg := config.NewConfig()
	fmt.Printf("%+v\n", cfg)

	// 2. Получаем цитату и автора
	zenAdapter := adapters.NewZenQuotesAdapter(cfg.QuoteAPIURL)
	quoteFetcher := usecases.NewQuoteFetcher(zenAdapter)

	quote, err := quoteFetcher.FetchFormattedQuote()
	if err != nil {
		fmt.Println("Ошибка при получении цитаты:", err)
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", quote)

	// 3. Переводим цитату
	trAdapter := adapters.NewTranslateAdapter(cfg.TranslateAPIURL)
	translateFetcher := usecases.NewTranslateFetcher(trAdapter)

	translateQuote, err := translateFetcher.FetchTranslated(quote, "en", "ru")
	if err != nil {
		fmt.Println("Ошибка при переводе цитаты:", err)
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", translateQuote)

	// 4. Отправим данные в TelegramBot
	tgAdapter, err := adapters.NewTelegramAdapter(cfg.TelegramToken)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	tgFetcher := usecases.NewTelegramFetcher(tgAdapter)

	formatted := usecases.FormatQuoteWithEmoji(translateQuote.Text, translateQuote.Author)
	sendMsgError := tgFetcher.FetchTelegram(cfg.TelegramChatId, formatted)
	fmt.Printf("%+v\n", sendMsgError)

}
