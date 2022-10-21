package telegram

type IDExtractor interface {
	ExtractID(url string) (string, error)
}
