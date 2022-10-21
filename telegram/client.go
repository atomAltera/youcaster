package telegram

import (
	e "github.com/atomAltera/youcaster/entities"
	"github.com/atomAltera/youcaster/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ListenConf struct {
	RestrictToChatIDs []int64
}

type Telegram struct {
	log       logger.Logger
	extractor IDExtractor
	bot       *tgbotapi.BotAPI
}

func NewTelegramClient(l logger.Logger, token string, e IDExtractor) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Telegram{
		log:       l,
		extractor: e,
		bot:       bot,
	}, nil
}

func (t *Telegram) ListenRequests(conf ListenConf) <-chan e.Request {
	uc := tgbotapi.NewUpdate(0)
	uc.Timeout = 60

	updates := t.bot.GetUpdatesChan(uc)
	rc := make(chan e.Request, 10)

	go func() {
		for u := range updates {
			if u.Message == nil {
				t.log.Warn("Received update without message")
				continue
			}

			if u.Message.Chat == nil {
				t.log.Warn("Received update without chat")
				continue
			}

			if u.Message.From == nil {
				t.log.Warn("Received update without from")
				continue
			}

			l := t.log.WithFields(map[string]interface{}{
				"chat_id":       u.Message.Chat.ID,
				"message_id":    u.Message.MessageID,
				"from_id":       u.Message.From.ID,
				"from_username": u.Message.From.UserName,
				"text_length":   len(u.Message.Text),
			})

			if len(conf.RestrictToChatIDs) > 0 {
				var found bool
				for _, id := range conf.RestrictToChatIDs {
					if u.Message.Chat.ID == id {
						found = true
						break
					}
				}
				if !found {
					l.Warn("Received update from restricted chat")
					continue
				}
			}

			if u.Message == nil {
				l.Warn("Received update without message")
				continue
			}

			vid, err := t.extractor.ExtractID(u.Message.Text)
			if err != nil {
				l.WithError(err).Error("failed to extract video id")
				continue
			}

			var r = e.Request{
				YoutubeVideoID: vid,
				TgChatID:       u.Message.Chat.ID,
				TgMessageID:    u.Message.MessageID,
			}

			rc <- r
		}
	}()

	return rc
}
