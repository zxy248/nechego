package pictures

import (
	"context"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Hello struct {
	Queries *data.Queries
}

var helloRe = handlers.NewRegexp("^!(п[рл]ив[а-я]*|хай|зд[ао]ров[а-я]*|ку|здрав[а-я]*)")

func (h *Hello) Match(c tele.Context) bool {
	return helloRe.MatchString(c.Text())
}

func (h *Hello) Handle(c tele.Context) error {
	ctx := context.Background()
	sticker, err := h.Queries.RandomHelloSticker(ctx)
	if err != nil {
		return err
	}
	return c.Send(&sticker)
}
