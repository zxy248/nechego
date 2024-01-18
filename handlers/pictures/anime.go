package pictures

import (
	"fmt"
	"math/rand"
	"nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Anime struct{}

func (h *Anime) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!аним", "!мульт")
}

func (h *Anime) Handle(c tele.Context) error {
	const spec = "https://thisanimedoesnotexist.ai/results/psi-%s/seed%05d.png"
	psi := [...]string{"0.3", "0.4", "0.5", "0.6", "0.7", "0.8", "0.9", "1.0",
		"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "2.0"}
	r1 := psi[rand.Intn(len(psi))]
	r2 := rand.Intn(100000)
	url := fmt.Sprintf(spec, r1, r2)
	return c.Send(&tele.Photo{File: tele.FromURL(url)})
}
