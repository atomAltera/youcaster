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
	created := time.Now()
	var updated time.Time
	for _, r := range rs {
		if r.UpdatedAt.After(updated) {
			updated = r.UpdatedAt
		}

		if r.CreatedAt.Before(created) {
			created = r.CreatedAt
		}
	}

	feed := &feeds.Feed{
		Title:       fb.Title,
		Link:        &feeds.Link{Href: "https://youtube.com"}, // TODO: make configurable
		Description: fb.Description,
		Author:      &feeds.Author{Name: fb.AuthorName, Email: fb.AuthorEmail},
		Updated:     updated,
		Created:     created,
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

	for _, r := range rs {
		feed.Add(&feeds.Item{
			Title: r.VideoInfo.Title,
			Link: &feeds.Link{
				Href: "https://www.youtube.com/watch?v=" + r.YoutubeVideoID,
			},
			Source: &feeds.Link{
				Href: "https://www.youtube.com/watch?v=" + r.YoutubeVideoID,
			},
			Author:      nil,
			Description: r.VideoInfo.Description,
			Id:          r.ID,
			Updated:     r.UpdatedAt,
			Created:     r.CreatedAt,
			Enclosure: &feeds.Enclosure{
				Url:    fb.PublicBaseURL + fmt.Sprintf(fb.FilePathPattern, r.FileName),
				Length: fmt.Sprintf("%d", r.FileSize),
				Type:   "audio/mpeg",
			},
			Content: r.VideoInfo.Description,
		})
	}

	xml, err := feed.ToRss()
	if err != nil {
		return "", fmt.Errorf("failed to build feed: %w", err)
	}

	return xml, err
}
