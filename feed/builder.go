package feed

import (
	"fmt"
	e "github.com/atomAltera/youcaster/entities"
	"github.com/eduncan911/podcast"
	"time"
)

type Builder struct {
	Title       string
	Description string
	AuthorName  string
	AuthorEmail string
	Copyright   string

	PublicBaseURL string
	MainLogoPath  string

	FilePathPattern string
}

func (b *Builder) BuildFeed(rs []e.Request) (string, error) {
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

	p := podcast.New(
		b.Title,
		"https://youtube.com",
		b.Description,
		&created,
		&updated,
	)

	p.AddImage(b.PublicBaseURL + b.MainLogoPath)
	p.AddAuthor(b.AuthorName, b.AuthorEmail)

	for _, r := range rs {
		i := podcast.Item{
			//XMLName:            xml.Name{},
			GUID:        r.ID,
			Title:       r.VideoInfo.Title,
			Link:        "https://youtube.com/watch?v=" + r.YoutubeVideoID,
			Description: r.VideoInfo.Description,
			//Author:             nil,
			//AuthorFormatted:    "",
			//Category:           "",
			//Comments:           "",
			Source: "https://youtube.com/watch?v=" + r.YoutubeVideoID,
			//PubDate:            nil,
			//PubDateFormatted:   "",
			//Enclosure:          nil,
			//IAuthor:            "",
			//ISubtitle:          "",
			//ISummary:           nil,
			//IImage:             nil,
			//IDuration:          "",
			//IExplicit:          "",
			//IIsClosedCaptioned: "",
			//IOrder:             "",
		}

		i.AddPubDate(&r.CreatedAt)
		i.AddImage(r.VideoInfo.ThumbnailURL)
		i.AddSummary(r.VideoInfo.Description)
		i.AddEnclosure(b.PublicBaseURL+fmt.Sprintf(b.FilePathPattern, r.FileName), podcast.MP3, r.FileSize)
		if r.VideoInfo.Duration > 0 {
			i.AddDuration(int64(r.VideoInfo.Duration.Seconds()))
		}

		if _, err := p.AddItem(i); err != nil {
			return "", fmt.Errorf("failed to add item to podcast: %w", err)
		}
	}

	return p.String(), nil
}
