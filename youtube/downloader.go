package youtube

import (
	"context"
	"fmt"
	"github.com/atomAltera/youcaster/logger"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

const ytDlpTemplate = `yt-dlp --extract-audio --audio-format=mp3 --audio-quality=0 -f m4a/bestaudio "%s" --no-progress -o "%s"`

type Downloader struct {
	log logger.Logger
	dir string
}

func NewDownloader(l logger.Logger, dir string) (*Downloader, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of %s: %w", dir, err)
	}

	return &Downloader{
		log: l,
		dir: absDir,
	}, nil
}

func (d *Downloader) Download(ctx context.Context, url string, filename string) (int64, error) {
	dest := path.Join(d.dir, filename)

	script := fmt.Sprintf(ytDlpTemplate, url, dest)

	cmd := exec.CommandContext(ctx, "sh", "-c", script)
	//cmd.Stdin = os.Stdin
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = d.dir

	d.log.Infof("executing command: %s", script)

	if err := cmd.Run(); err != nil {
		_ = os.Remove(dest)

		return 0, fmt.Errorf("failed to execute command: %v", err)
	}

	s, err := os.Stat(dest)
	if err != nil {
		return 0, fmt.Errorf("failed to get file size of %s: %w", dest, err)
	}

	return s.Size(), nil
}
