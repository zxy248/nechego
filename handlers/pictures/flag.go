package pictures

import (
	"fmt"
	"math/rand"
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Flag struct{}

func (h *Flag) Match(s string) bool {
	return handlers.HasPrefix(s, "!флаг")
}

func (h *Flag) Handle(c tele.Context) error {
	return c.Send(&tele.Photo{File: tele.FromURL(randomFlagURL())})
}

func randomFlagURL() string {
	const format = "https://thisflagdoesnotexist.com/images/%d.png"
	return fmt.Sprintf(format, rand.Intn(5000))
}
