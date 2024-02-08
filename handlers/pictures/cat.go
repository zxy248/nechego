package pictures

import (
	"fmt"
	"math/rand/v2"
	"net/http"

	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Cat struct{}

func (h *Cat) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!кот", "!кош")
}

func (h *Cat) Handle(c tele.Context) error {
	const spec = "https://d2ph5fj80uercy.cloudfront.net/%02d/cat%d.jpg"
	url := fmt.Sprintf(spec, 1+rand.N(6), rand.N(5000))
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(resp.Body)})
}
