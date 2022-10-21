package feed

import (
	"fmt"
	e "github.com/atomAltera/youcaster/entities"
	"github.com/gorilla/feeds"
	"time"
)

type Builder struct {
	Title       string
	Description string
	AuthorName  string
	AuthorEmail string
	Copyright   string

	PublicBaseURL  string
	MainLogoPath   string
	MainLogoWidth  int
	MainLogoHeight int

	FilePathPattern string
}

func (fb *Builder) BuildFeed(rs []e.Request) (string, error) {
	feed := &feeds.Feed{
		Title:       fb.Title,
		Link:        &feeds.Link{Href: "https://youtube.com"}, // TODO: make configurable
		Description: fb.Description,
		Author:      &feeds.Author{Name: fb.AuthorName, Email: fb.AuthorEmail},
		Updated:     time.Time{}, // TODO: Provide this
		Created:     time.Date(2022, 10, 15, 0, 0, 0, 0, time.UTC),
		Id:          "",
		Subtitle:    "",
		Items:       nil,
		Copyright:   fb.Copyright,
		Image: &feeds.Image{
			Title:  fb.Title,
			Url:    fb.PublicBaseURL + fb.MainLogoPath,
			Link:   fb.PublicBaseURL + fb.MainLogoPath,
			Width:  fb.MainLogoWidth,
			Height: fb.MainLogoHeight,
		},
	}

	feed.Items = make([]*feeds.Item, len(rs))
	for i, r := range rs {
		feed.Items[i] = &feeds.Item{
			Title: r.VideoInfo.Title,
			Link: &feeds.Link{
				Href: "https://www.youtube.com/watch?v=" + r.YoutubeVideoID,
			},
			Source:      nil,
			Author:      nil,
			Description: r.VideoInfo.Description,
			Id:          r.ID,
			Updated:     time.Time{},
			Created:     r.CreatedAt,
			Enclosure: &feeds.Enclosure{
				Url:    fb.PublicBaseURL + fmt.Sprintf(fb.FilePathPattern, r.FileName),
				Length: fmt.Sprintf("%d", r.FileSize),
				Type:   "audio/mpeg",
			},
			Content: r.VideoInfo.Description,
		}
	}

	xml, err := feed.ToRss()
	if err != nil {
		return "", fmt.Errorf("failed to build feed: %w", err)
	}

	return xml, err
}
