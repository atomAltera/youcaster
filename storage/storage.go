package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	var err error
	if basePath, err = filepath.Abs(basePath); err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	if err = os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	return &FileStorage{
		basePath: basePath,
	}, nil
}

func (s *FileStorage) ReadFile(ctx context.Context, name string) (io.ReadCloser, error) {
	fp := path.Join(s.basePath, name)
	f, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return f, nil
}
