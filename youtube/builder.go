package youtube

type URLBuilder struct {
}

func NewURLBuilder() *URLBuilder {
	return &URLBuilder{}
}

func (b *URLBuilder) BuildVideoURL(id string) string {
	return "https://youtu.be/" + id
}
