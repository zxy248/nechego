package fun

import (
	"context"
	"html"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"
	tu "github.com/zxy248/nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Who struct {
	Queries *data.Queries
}

var whoRe = handlers.NewRegexp("^!кто(.*)")

func (h *Who) Match(c tele.Context) bool {
	return whoRe.MatchString(c.Text())
}

func (h *Who) Handle(c tele.Context) error {
	ctx := context.Background()
	arg := data.RandomUsersParams{
		ChatID: c.Chat().ID,
		Limit:  1,
	}
	users, err := h.Queries.RandomUsers(ctx, arg)
	if err != nil {
		return err
	}

	out := tu.Link(c, users[0].ID) + getWhoArgument(c.Text())
	return c.Send(out, tele.ModeHTML)
}

func getWhoArgument(s string) string {
	return html.EscapeString(whoRe.FindStringSubmatch(s)[1])
}
