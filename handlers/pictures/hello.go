package pictures

import (
	"context"
	"encoding/json"

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
	raw, err := h.Queries.RandomHelloSticker(ctx)
	if err != nil {
		return err
	}

	var s tele.Sticker
	if err := json.Unmarshal(raw, &s); err != nil {
		return err
	}
	return c.Send(&s)
}
