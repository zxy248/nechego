package pictures

import (
	"fmt"
	"math/rand"
	"nechego/handlers"
	"net/http"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Cat struct{}

func (h *Cat) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!кот", "!кош")
}

func (h *Cat) Handle(c tele.Context) error {
	const spec = "https://d2ph5fj80uercy.cloudfront.net/%02d/cat%d.jpg"
	r1 := 1 + rand.Intn(6)
	r2 := rand.Intn(5000)
	url := fmt.Sprintf(spec, r1, r2)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(resp.Body)})
}
