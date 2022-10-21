package web

import (
	"context"
	e "github.com/atomAltera/youcaster/entities"
	"io"
)

type FileReader interface {
	ReadFile(ctx context.Context, name string) (io.ReadCloser, error)
}

type RequestsStore interface {
	List(ctx context.Context, ss []e.RequestStatus) ([]e.Request, error)
}

type FeedBuilder interface {
	BuildFeed(rs []e.Request) (string, error)
}
