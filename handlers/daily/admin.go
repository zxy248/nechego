package daily

import (
	"context"
	"fmt"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"
	tu "github.com/zxy248/nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Admin struct {
	Queries *data.Queries
}

var adminRe = handlers.NewRegexp("^!админ")

func (h *Admin) Match(c tele.Context) bool {
	return adminRe.MatchString(c.Text())
}

func (h *Admin) Handle(c tele.Context) error {
	ctx := context.Background()
	chat, err := h.Queries.GetChat(ctx, c.Chat().ID)
	if err != nil {
		return err
	}
	link := tu.Link(c, chat.Data.Admin)
	out := fmt.Sprintf("<b>Админ дня</b> — %s 👑", link)
	return c.Send(out, tele.ModeHTML)
}
