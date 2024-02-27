package handlers

import (
	"strings"

	"github.com/zxy248/nechego/game"
	tu "github.com/zxy248/nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Pass struct {
	Universe *game.Universe
	Logger   *Logger
}

func (h *Pass) Match(c tele.Context) bool {
	return true
}

func (h *Pass) Handle(c tele.Context) error {
	h.saveSticker(c)
	return h.logMessage(c)
}

func (h *Pass) logMessage(c tele.Context) error {
	text := strings.TrimSpace(c.Text())
	if !strings.ContainsRune(text, '\n') && text != "" && len(text) < 1024 {
		if err := h.Logger.Log(c.Chat().ID, text); err != nil {
			return err
		}
	}
	return nil
}

func (h *Pass) saveSticker(c tele.Context) {
	const limit = 200

	if s := c.Message().Sticker; s != nil {
		w := tu.Lock(c, h.Universe)
		defer w.Unlock()

		w.Stickers = append(w.Stickers, s.FileID)
		if len(w.Stickers) > limit {
			w.Stickers = w.Stickers[len(w.Stickers)-limit:]
		}
	}
}
