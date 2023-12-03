package pictures

import (
	"fmt"
	"math/rand"
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Flag struct{}

func (h *Flag) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!флаг")
}

func (h *Flag) Handle(c tele.Context) error {
	const spec = "https://thisflagdoesnotexist.com/images/%d.png"
	r := rand.Intn(5000)
	url := fmt.Sprintf(spec, r)
	return c.Send(&tele.Photo{File: tele.FromURL(url)})
}
