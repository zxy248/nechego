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
	users, err := h.Queries.RecentUsers(ctx, c.Chat().ID)
	if err != nil {
		return err
	}
	list := randomSample(users, 3+rand.N(3))
	name := topName(c.Text())
	out := fmt.Sprintf("<b>🏆 Топ %s</b>\n", name)
	for i, u := range list {
		l := tu.Link(c, u)
		out += fmt.Sprintf("%d. <b>%s</b>\n", i+1, l)
	}
	return c.Send(out, tele.ModeHTML)
}

func topName(s string) string {
	return topRe.FindStringSubmatch(s)[1]
}
