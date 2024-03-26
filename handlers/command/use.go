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
	if _, err := h.Queries.SelectCommand(context.Background(), data.SelectCommandParams{
		ChatID:     c.Chat().ID,
		Definition: commandDefinition(c.Text()),
	}); err != nil {
		return false
	}
	return true
}

func (h *Use) Handle(c tele.Context) error {
	cmd, err := h.Queries.SelectCommand(context.Background(), data.SelectCommandParams{
		ChatID:     c.Chat().ID,
		Definition: commandDefinition(c.Text()),
	})
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
