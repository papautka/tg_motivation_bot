package interfaces

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

type TelegramApi interface {
	SendMessage(chatId int, message string) error
}
