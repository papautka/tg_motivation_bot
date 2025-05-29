package usecases

import (
	"tg_motivation_bot/internal/interfaces"
)

type QuoteFetcher struct {
	Provider interfaces.QuoteProvider
}

func NewQuoteFetcher(p interfaces.QuoteProvider) *QuoteFetcher {
	return &QuoteFetcher{Provider: p}
}

func (qf *QuoteFetcher) FetchFormattedQuote() (*interfaces.Quote, error) {
	quote, err := qf.Provider.GetQuote()
	if err != nil {
		return nil, err
	}
	return quote, nil
}

func FormatQuoteWithEmoji(text, author string) string {
	return "💬 \"" + text + "\"\n🖊️ — " + author
}
