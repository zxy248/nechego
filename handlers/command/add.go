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

	if err := h.Queries.AddCommand(context.Background(), data.AddCommandParams{
		ChatID:            c.Chat().ID,
		Definition:        commandDefinition(match[1]),
		SubstitutionText:  commandSubstitution(match[2]),
		SubstitutionPhoto: photo,
	}); err != nil {
		return err
	}
	return c.Send("✅ Команда добавлена.")
}
