package pictures

import (
	"nechego/handlers"
	"net/http"

	tele "gopkg.in/telebot.v3"
)

type Soy struct{}

func (h *Soy) Match(s string) bool {
	return handlers.HasPrefix(s, "!сой")
}

func (h *Soy) Handle(c tele.Context) error {
	resp, err := http.Get("https://booru.soy/random_image/download")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(resp.Body)})
}
