package usecases

import "tg_motivation_bot/internal/interfaces"

type TranslateFetcher struct {
	Provider interfaces.Translator
}

func NewTranslateFetcher(p interfaces.Translator) *TranslateFetcher {
	return &TranslateFetcher{
		Provider: p,
	}
}

func (t *TranslateFetcher) FetchTranslated(quote *interfaces.Quote, fromLang, toLang string) (*interfaces.Quote, error) {
	translate, err := t.Provider.Translate(quote, fromLang, toLang)
	if err != nil {
		return nil, err
	}
	return translate, nil
}
