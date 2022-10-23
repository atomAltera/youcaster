package worker

import (
	"context"
	e "github.com/atomAltera/youcaster/entities"
)

type InfoGetter interface {
	GetInfo(ctx context.Context, id string) (*e.VideoInfo, error)
}

type Downloader interface {
	Download(ctx context.Context, id string, filename string) (int64, error)
}

type RequestsStore interface {
	Create(ctx context.Context, r e.Request) error
	Update(ctx context.Context, r e.Request) error
	List(ctx context.Context, ss []e.RequestStatus) ([]e.Request, error)
}
