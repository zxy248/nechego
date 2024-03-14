package daily

import (
	"context"
	"fmt"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"
	tu "github.com/zxy248/nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Eblan struct {
	Queries *data.Queries
}

var eblanRe = handlers.NewRegexp("^![ие][б6п]?л[ап]н[а-я]*")

func (h *Eblan) Match(c tele.Context) bool {
	return eblanRe.MatchString(c.Text())
}

func (h *Eblan) Handle(c tele.Context) error {
	ctx := context.Background()
	chat, err := h.Queries.GetChat(ctx, c.Chat().ID)
	if err != nil {
		return err
	}
	link := tu.Link(c, chat.Data.Eblan)
	out := fmt.Sprintf("<b>Еблан дня</b> — %s 😸", link)
	return c.Send(out, tele.ModeHTML)
}
