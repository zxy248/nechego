package pictures

import (
	"fmt"
	"math/rand"
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Anime struct{}

func (h *Anime) Match(s string) bool {
	return handlers.MatchRegexp("^!(аним|мульт)", s)
}

func (h *Anime) Handle(c tele.Context) error {
	return c.Send(&tele.Photo{File: tele.FromURL(randomAnimeURL())})
}

func randomAnimeURL() string {
	const format = "https://thisanimedoesnotexist.ai/results/psi-%s/seed%05d.png"
	psis := [...]string{"0.3", "0.4", "0.5", "0.6", "0.7", "0.8", "0.9", "1.0",
		"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "2.0"}
	psi := psis[rand.Intn(len(psis))]
	return fmt.Sprintf(format, psi, rand.Intn(100_000))
}
