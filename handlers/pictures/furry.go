package pictures

import (
	"fmt"
	"math/rand/v2"

	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Furry struct{}

func (h *Furry) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!фур")
}

func (h *Furry) Handle(c tele.Context) error {
	const spec = "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%05d.jpg"
	url := fmt.Sprintf(spec, rand.N(100000))
	return c.Send(&tele.Photo{File: tele.FromURL(url)})
}
