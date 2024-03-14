package fun

import (
	"context"
	"fmt"
	"math/rand/v2"
	"slices"

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
	users, err := h.Queries.RecentUsers(ctx, c.Chat().ID)
	if err != nil {
		return err
	}
	list := randomSample(users, 3+rand.N(3))
	name := listName(c.Text())
	out := fmt.Sprintf("<b>📝 Список %s</b>\n", name)
	for _, u := range list {
		l := tu.Link(c, u)
		out += fmt.Sprintf("• <b>%s</b>\n", l)
	}
	return c.Send(out, tele.ModeHTML)
}

func listName(s string) string {
	return listRe.FindStringSubmatch(s)[1]
}

func randomSample[T any](a []T, n int) []T {
	c := slices.Clone(a)
	rand.Shuffle(len(c), func(i, j int) {
		c[i], c[j] = c[j], c[i]
	})
	return c[:min(len(c), n)]
}
