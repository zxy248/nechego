package command

import (
	"context"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Remove struct {
	Queries *data.Queries
}

var removeRe = handlers.NewRegexp("^!(удалить|убрать) (" + definitionPattern + ")")

func (h *Remove) Match(c tele.Context) bool {
	return removeRe.MatchString(c.Text())
}

func (h *Remove) Handle(c tele.Context) error {
	match := removeRe.FindStringSubmatch(c.Text())
	ctx := context.Background()
	arg := data.DeleteCommandsParams{
		ChatID:     c.Chat().ID,
		Definition: sanitizeDefinition(match[2]),
	}
	if err := h.Queries.DeleteCommands(ctx, arg); err != nil {
		return err
	}
	return c.Send("❌ Команда удалена.")
}
