package fun

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"
	tu "github.com/zxy248/nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Top struct {
	Queries *data.Queries
}

var topRe = handlers.NewRegexp("^!топ ?(.*)")

func (h *Top) Match(c tele.Context) bool {
	return topRe.MatchString(c.Text())
}

func (h *Top) Handle(c tele.Context) error {
	ctx := context.Background()
	arg := data.RandomUsersParams{
		ChatID: c.Chat().ID,
		Limit:  3 + rand.N[int32](3),
	}
	users, err := h.Queries.RandomUsers(ctx, arg)
	if err != nil {
		return err
	}

	out := fmt.Sprintf("<b>🏆 Топ %s</b>\n", getTopArgument(c.Text()))
	for i, u := range users {
		out += fmt.Sprintf("%d. <b>%s</b>\n", i+1, tu.Link(c, u))
	}
	return c.Send(out, tele.ModeHTML)
}

func getTopArgument(s string) string {
	return topRe.FindStringSubmatch(s)[1]
}
