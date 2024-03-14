package command

import (
	"context"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Add struct {
	Queries *data.Queries
}

var addRe = handlers.NewRegexp("^!добавить (" + definitionPattern + ")\\|?(.*)")

func (h *Add) Match(c tele.Context) bool {
	return addRe.MatchString(c.Text())
}

func (h *Add) Handle(c tele.Context) error {
	match := addRe.FindStringSubmatch(c.Text())
	var photo string
	if p := c.Message().Photo; p != nil {
		photo = p.FileID
	}
	ctx := context.Background()
	arg := data.AddCommandParams{
		ChatID:            c.Chat().ID,
		Definition:        sanitizeDefinition(match[1]),
		SubstitutionText:  sanitizeSubstitution(match[2]),
		SubstitutionPhoto: photo,
	}
	if err := h.Queries.AddCommand(ctx, arg); err != nil {
		return err
	}
	return c.Send("✅ Команда добавлена.")
}
