package pictures

import (
	"encoding/json"
	"math/rand"
	"nechego/handlers"
	"os"

	tele "gopkg.in/telebot.v3"
)

type Hello struct {
	Path  string
	cache []tele.Sticker
}

func (h *Hello) Match(s string) bool {
	return handlers.MatchRegexp("^!(п[рл]ив[а-я]*|хай|зд[ао]ров[а-я]*|ку|здрав[а-я]*)", s)
}

func (h *Hello) Handle(c tele.Context) error {
	if err := h.init(); err != nil {
		return err
	}
	return c.Send(&h.cache[rand.Intn(len(h.cache))])
}

func (h *Hello) init() error {
	if h.cache == nil {
		ss, err := loadStickers(h.Path)
		if err != nil {
			return err
		}
		h.cache = ss
	}
	return nil
}

func loadStickers(path string) ([]tele.Sticker, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := []tele.Sticker{}
	if err := json.NewDecoder(f).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}
