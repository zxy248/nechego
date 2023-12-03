package pictures

import (
	"nechego/handlers"
	"net/http"

	tele "gopkg.in/telebot.v3"
)

type Soy struct{}

func (h *Soy) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!сой")
}

func (h *Soy) Handle(c tele.Context) error {
	const url = "https://booru.soy/random_image/download"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(resp.Body)})
}
