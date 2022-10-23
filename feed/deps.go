package feed

type URLBuilder interface {
	BuildVideoURL(id string) string
}
