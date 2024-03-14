package fun

import (
	"context"
	"html"
	"math/rand/v2"

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
	users, err := h.Queries.RecentUsers(ctx, c.Chat().ID)
	if err != nil {
		return err
	}
	u := users[rand.N(len(users))]
	w := parseWho(c.Text())
	l := tu.Link(c, u.ID)
	s := l + w
	return c.Send(s, tele.ModeHTML)
}

func parseWho(s string) string {
	return html.EscapeString(whoRe.FindStringSubmatch(s)[1])
}
