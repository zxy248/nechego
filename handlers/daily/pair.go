package daily

import (
	"context"
	"fmt"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"
	tu "github.com/zxy248/nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Pair struct {
	Queries *data.Queries
}

var pairRe = handlers.NewRegexp("^!пара")

func (h *Pair) Match(c tele.Context) bool {
	return pairRe.MatchString(c.Text())
}

func (h *Pair) Handle(c tele.Context) error {
	ctx := context.Background()
	chat, err := h.Queries.GetChat(ctx, c.Chat().ID)
	if err != nil {
		return err
	}
	link1 := tu.Link(c, chat.Data.Pair1)
	link2 := tu.Link(c, chat.Data.Pair2)
	out := fmt.Sprintf("<b>✨ Пара дня</b> — %s 💘 %s", link1, link2)
	return c.Send(out, tele.ModeHTML)
}
