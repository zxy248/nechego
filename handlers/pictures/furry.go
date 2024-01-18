package pictures

import (
	"fmt"
	"github.com/zxy248/nechego/handlers"
	"math/rand"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Furry struct{}

func (h *Furry) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!фур")
}

func (h *Furry) Handle(c tele.Context) error {
	const spec = "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%05d.jpg"
	r := rand.Intn(100000)
	url := fmt.Sprintf(spec, r)
	return c.Send(&tele.Photo{File: tele.FromURL(url)})
}
