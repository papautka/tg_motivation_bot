package adapters

type ResponseData struct {
	TranslatedText string `json:"translatedText"`
}

type MyMemoryResponse struct {
	ResponseData ResponseData `json:"responseData"`
}
