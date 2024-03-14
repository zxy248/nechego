package handlers

import (
	"context"

	"github.com/zxy248/nechego/data"

	tele "gopkg.in/zxy248/telebot.v3"
)

const MessageIDKey = "database_message_id"

type Pass struct {
	Queries *data.Queries
}

func (h *Pass) Match(c tele.Context) bool {
	return true
}

func (h *Pass) Handle(c tele.Context) error {
	ctx := context.Background()
	messageID := c.Get(MessageIDKey).(int64)
	if err := h.Queries.SetMessageNotCommand(ctx, messageID); err != nil {
		return err
	}
	if s := c.Message().Sticker; s != nil {
		arg := data.AddStickerParams{
			MessageID: messageID,
			FileID:    s.FileID,
		}
		if err := h.Queries.AddSticker(ctx, arg); err != nil {
			return err
		}
	}
	return nil
}
