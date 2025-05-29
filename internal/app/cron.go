package app

import (
	"log/slog"
	"tg_motivation_bot/internal/config"
	"tg_motivation_bot/internal/sheduler"
	"tg_motivation_bot/internal/usecases"
)

func startCron(cfg *config.Config, quoteFetcher *usecases.QuoteFetcher, translateFetcher *usecases.TranslateFetcher, telegramFetcher *usecases.TelegramFetcher) {
	cronExpr := "* * * * *" // Каждую минуту
	sheduler.InitScheduler(cronExpr, func() {
		// 4.0 получаем цитату на англосаксонском
		msg, err := getTranslatedQuote(quoteFetcher, translateFetcher)
		if err != nil {
			slog.Error("Ошибка при получении или переводе", slog.String("error", err.Error()))
			return
		}
		// 4.4 отправляем в tg
		// в startCron мы не передаем клавиатуру поэтому в параметры клавиатуры передаем nil
		errSendMsg := telegramFetcher.FetchTelegram(cfg.TelegramChatId, msg, nil)
		if errSendMsg != nil {
			slog.Error("Ошибка при отправке цитаты в telegram", slog.String("error", errSendMsg.Error()))
			return
		} else {
			slog.Info("Цитата успешно отправлена")
		}
	})
}

func getTranslatedQuote(quoteFetcher *usecases.QuoteFetcher, translateFetcher *usecases.TranslateFetcher) (string, error) {
	// 4.1 получаем цитату на англосаксонском
	myQuote, errGetQuote := quoteFetcher.FetchFormattedQuote()
	if errGetQuote != nil {
		slog.Error("Ошибка при получении цитаты", slog.String("error", errGetQuote.Error()))
		return "", errGetQuote
	}
	// 4.2 переводим на великий могучий
	mySlavicQuote, errTranslate := translateFetcher.FetchTranslated(myQuote, "en", "ru")
	if errTranslate != nil {
		slog.Error("Ошибка при переводе цитаты", slog.String("error", errTranslate.Error()))
		return "", errTranslate
	}

	// 4.3 форматируем строку
	msg := usecases.FormatQuoteWithEmoji(mySlavicQuote.Text, mySlavicQuote.Author)
	return msg, nil
}
