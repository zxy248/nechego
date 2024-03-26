package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

	cmd, err := h.Queries.DeleteCommand(context.Background(), data.DeleteCommandParams{
		ChatID:     c.Chat().ID,
		Definition: commandDefinition(match[2]),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return c.Send("⚠️ Такой команды нет.")
	}
	if err != nil {
		return err
	}

	if cmd.SubstitutionText != "" {
		const format = "❎ Команда удалена.\n\n" +
			"<i>Чтобы вернуть команду, используйте <code>!добавить %s|%s</code></i>"
		return c.Send(fmt.Sprintf(format, cmd.Definition, cmd.SubstitutionText), tele.ModeHTML)
	}
	return c.Send("❎ Команда удалена.")
}
