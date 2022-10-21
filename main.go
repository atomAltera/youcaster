package main

import (
	"github.com/atomAltera/youcaster/database"
	"github.com/atomAltera/youcaster/logger"
	"github.com/atomAltera/youcaster/telegram"
	"github.com/atomAltera/youcaster/worker"
	"github.com/atomAltera/youcaster/youtube"
	"github.com/umputun/go-flags"
	"os"
)

var opts struct {
	Log logger.Opts `group:"Logging Options" namespace:"log" env-namespace:"LOG"`

	PublicBaseURL string `long:"public-base-url" env:"PUBLIC_BASE_URL" description:"public base url of the server"`
	DownloadPath  string `long:"download-path" env:"DOWNLOAD_PATH" description:"path to download videos to" default:"/tmp"`

	Mongo struct {
		URI string `long:"uri" env:"URI" required:"true" description:"Mongodb database uri"`
	} `group:"MongoDB Options" namespace:"mongo" env-namespace:"MONGO"`

	Web struct {
		Addr string `long:"addr" env:"ADDR" description:"address to listen on" default:"0.0.0.0:3000"`
	} `group:"Web Options" namespace:"web" env-namespace:"WEB"`

	Telegram struct {
		Token   string  `long:"token" env:"TOKEN" required:"true" description:"Telegram bot token"`
		ChatIDs []int64 `short:"c" long:"chat" env:"CHATS" description:"Telegram chat ids"`
	} `group:"Telegram Options" namespace:"telegram" env-namespace:"TELEGRAM"`

	Google struct {
		APIKey string `long:"api-key" env:"API_KEY" required:"true" description:"Google API key"`
	} `group:"Google Options" namespace:"google" env-namespace:"GOOGLE"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	log := logger.GetLogger(opts.Log)

	db, err := database.New(opts.Mongo.URI)
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	yig, err := youtube.NewInfoGetter(
		opts.Google.APIKey,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to create youtube info getter: %v", err)
	}

	ydl, err := youtube.NewDownloader(
		log.WithField("module", "downloader"),
		opts.DownloadPath,
	)
	if err != nil {
		log.Fatalf("failed to create youtube downloader: %v", err)
	}

	w := worker.NewWorker(
		log.WithField("module", "worker"),
		db.Requests,
		yig,
		ydl,
	)

	tg, err := telegram.NewTelegramClient(
		log.WithField("module", "telegram"),
		opts.Telegram.Token,
		youtube.NewURLParser(),
	)
	if err != nil {
		log.Fatalf("failed to create telegram client: %v", err)
	}

	rc := tg.ListenRequests(telegram.ListenConf{
		RestrictToChatIDs: opts.Telegram.ChatIDs,
	})

	w.StartListenRequests(rc)

	w.StartProcessingRequests()

	// Create worker
}
