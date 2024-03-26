package command

import (
	"context"

	"github.com/zxy248/nechego/data"
	tele "gopkg.in/zxy248/telebot.v3"
)

type Use struct {
	Queries *data.Queries
}

func (h *Use) Match(c tele.Context) bool {
	ctx := context.Background()
	arg := data.SelectCommandParams{
		ChatID:     c.Chat().ID,
		Definition: c.Text(),
	}
	_, err := h.Queries.SelectCommand(ctx, arg)
	if err != nil {
		return false
	}
	return true
}

func (h *Use) Handle(c tele.Context) error {
	ctx := context.Background()
	arg := data.SelectCommandParams{
		ChatID:     c.Chat().ID,
		Definition: c.Text(),
	}
	cmd, err := h.Queries.SelectCommand(ctx, arg)
	if err != nil {
		return err
	}
	if cmd.SubstitutionPhoto != "" {
		f := tele.File{FileID: cmd.SubstitutionPhoto}
		p := &tele.Photo{File: f, Caption: cmd.SubstitutionText}
		return c.Send(p)
	}
	return c.Send(cmd.SubstitutionText)
}
