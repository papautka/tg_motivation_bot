package interfaces

type Quote struct {
	Text   string
	Author string
}

type QuoteProvider interface {
	GetQuote() (*Quote, error)
}
