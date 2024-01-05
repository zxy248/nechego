package pictures

import (
	"encoding/json"
	"math/rand"
	"nechego/handlers"
	"os"
	"sync"

	tele "gopkg.in/telebot.v3"
)

type Hello struct {
	Path string

	s  []tele.Sticker
	mu sync.Mutex
}

var helloRe = handlers.NewRegexp("^!(п[рл]ив[а-я]*|хай|зд[ао]ров[а-я]*|ку|здрав[а-я]*)")

func (h *Hello) Match(c tele.Context) bool {
	return helloRe.MatchString(c.Text())
}

func (h *Hello) Handle(c tele.Context) error {
	if err := h.init(); err != nil {
		return err
	}
	return c.Send(&h.s[rand.Intn(len(h.s))])
}

func (h *Hello) init() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.s == nil {
		s, err := loadStickers(h.Path)
		if err != nil {
			return err
		}
		h.s = s
	}
	return nil
}

func loadStickers(path string) ([]tele.Sticker, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s []tele.Sticker
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	return s, nil
}
