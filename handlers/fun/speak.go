package fun

import (
	"strings"
	"time"

	"github.com/zxy248/nechego/game"
	"github.com/zxy248/nechego/handlers"
	"github.com/zxy248/nechego/markov"
	tu "github.com/zxy248/nechego/teleutil"
	tele "gopkg.in/zxy248/telebot.v3"
)

type Speak struct {
	Universe *game.Universe
	Logger   *handlers.Logger
	Attempts int
}

func (h *Speak) Match(c tele.Context) bool {
	_, ok := parseSpeak(c.Text())
	return ok
}

func (h *Speak) Handle(c tele.Context) error {
	world := tu.Lock(c, h.Universe)
	defer world.Unlock()

	size, _ := parseSpeak(c.Text())
	chain, err := h.getChain(world)
	if err != nil {
		return err
	}
	res := h.generateText(chain, size)
	if res == "" {
		res = "⚠️ Не удалось сгенерировать сообщение."
	}
	return c.Send(res)
}

func (h *Speak) getChain(w *game.World) (*markov.Chain, error) {
	if time.Since(w.ChainUpdate) > 30*time.Minute {
		samples, err := h.Logger.Messages(w.ID)
		if err != nil {
			return nil, err
		}
		w.Chain = markov.New(samples)
	}
	return w.Chain, nil
}

func (h *Speak) generateText(c *markov.Chain, size int) string {
	minSize := [...]int{0, 2, 4, 8}
	maxSize := [...]int{100, 3, 7, 100}

	for range h.Attempts {
		s := c.Generate()
		if len(s) >= minSize[size] && len(s) <= maxSize[size] {
			return strings.Join(s, " ")
		}
	}
	return ""
}

func parseSpeak(s string) (size int, ok bool) {
	if len(s) > 50 {
		return 0, false
	}
	s = strings.ToLower(s)
	tokens := strings.Fields(s)
	if len(tokens) < 2 {
		return 0, false
	}
	w := tokens[0] == "w" || tokens[0] == "witless"
	g := tokens[1] == "s" || tokens[1] == "speak" || tokens[1] == "g" || tokens[1] == "generate"
	if w && g {
		if len(tokens) == 2 {
			return 0, true
		}
		return convertSize(tokens[2])
	}
	return 0, false
}

func convertSize(s string) (size int, ok bool) {
	switch s {
	case "any", "любое":
		return 0, true
	case "sm", "small", "маленькое", "короткое":
		return 1, true
	case "md", "medium", "среднее":
		return 2, true
	case "lg", "large", "большое", "длинное":
		return 3, true
	}
	return 0, false
}
