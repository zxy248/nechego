package handlers

import (
	"strings"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Pass struct {
	Logger *Logger
}

func (h *Pass) Match(c tele.Context) bool {
	return true
}

func (h *Pass) Handle(c tele.Context) error {
	text := strings.TrimSpace(c.Text())
	if !strings.ContainsRune(text, '\n') && text != "" && len(text) < 1024 {
		if err := h.Logger.Log(c.Chat().ID, text); err != nil {
			return err
		}
	}
	return nil
}
