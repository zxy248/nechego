package pictures

import (
	"fmt"
	"math/rand/v2"

	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Flag struct{}

func (h *Flag) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!флаг")
}

func (h *Flag) Handle(c tele.Context) error {
	const spec = "https://thisflagdoesnotexist.com/images/%d.png"
	url := fmt.Sprintf(spec, rand.N(5000))
	return c.Send(&tele.Photo{File: tele.FromURL(url)})
}
