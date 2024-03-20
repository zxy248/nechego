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

type List struct {
	Queries *data.Queries
}

var listRe = handlers.NewRegexp("^!список ?(.*)")

func (h *List) Match(c tele.Context) bool {
	return listRe.MatchString(c.Text())
}

func (h *List) Handle(c tele.Context) error {
	ctx := context.Background()
	arg := data.RandomUsersParams{
		ChatID: c.Chat().ID,
		Limit:  3 + rand.N[int32](3),
	}
	users, err := h.Queries.RandomUsers(ctx, arg)
	if err != nil {
		return err
	}

	out := fmt.Sprintf("<b>📝 Список %s</b>\n", getListArgument(c.Text()))
	for _, u := range users {
		out += fmt.Sprintf("• <b>%s</b>\n", tu.Link(c, u))
	}
	return c.Send(out, tele.ModeHTML)
}

func getListArgument(s string) string {
	return listRe.FindStringSubmatch(s)[1]
}
