package youtube

import (
	"errors"
	"regexp"
)

type URLParser struct {
	res []*regexp.Regexp
}

func NewURLParser() *URLParser {
	return &URLParser{
		res: []*regexp.Regexp{
			regexp.MustCompile(`https://(www.?)youtube.com/watch\?v=(?P<id>[a-zA-Z0-9_-]+)`),
			regexp.MustCompile(`https://youtu.be/(?P<id>[a-zA-Z0-9_-]+)`),
		},
	}
}

func (p *URLParser) Parse(url string) (string, error) {
	for _, re := range p.res {
		if re.MatchString(url) {
			matches := re.FindStringSubmatch(url)
			for i, name := range re.SubexpNames() {
				if name == "id" {
					return matches[i], nil
				}
			}
		}
	}

	return "", errors.New("invalid url")
}
