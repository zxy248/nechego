package pictures

import (
	"fmt"
	"math/rand"
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Furry struct{}

func (h *Furry) Match(s string) bool {
	return handlers.MatchPrefix("!фур", s)
}

func (h *Furry) Handle(c tele.Context) error {
	return c.Send(&tele.Photo{File: tele.FromURL(randomFurryURL())})
}

func randomFurryURL() string {
	const format = "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%05d.jpg"
	return fmt.Sprintf(format, rand.Intn(100_000))
}
