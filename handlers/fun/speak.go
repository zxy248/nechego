package fun

import (
	"context"
	"strings"
	"time"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"
	"github.com/zxy248/nechego/handlers/fun/markov"
	tele "gopkg.in/zxy248/telebot.v3"
)

type cachedChain struct {
	chain *markov.Chain
	time  time.Time
}

type Speak struct {
	Queries  *data.Queries
	Attempts int

	data handlers.Store[*cachedChain]
}

func (h *Speak) Match(c tele.Context) bool {
	_, ok := parseSpeak(c.Text())
	return ok
}

func (h *Speak) Handle(c tele.Context) error {
	v, done := h.data.Get(c.Chat().ID, &cachedChain{})
	if time.Since(v.time) > 10*time.Minute {
		ch, err := h.buildChain(c.Chat().ID)
		if err != nil {
			return err
		}
		v.chain = ch
		v.time = time.Now()
	}
	chain := v.chain
	done()

	size, _ := parseSpeak(c.Text())
	out := h.generateText(chain, size)
	if out == "" {
		out = "⚠️ Не удалось сгенерировать сообщение."
	}
	return c.Send(out)
}

func (h *Speak) buildChain(id int64) (*markov.Chain, error) {
	ctx := context.Background()
	messages, err := h.Queries.ListMessages(ctx, id)
	if err != nil {
		return nil, err
	}
	var samples []string
	for _, m := range messages {
		text := strings.ToLower(m.Content)
		samples = append(samples, text)
	}
	return markov.New(samples), nil
}

func (h *Speak) generateText(c *markov.Chain, size int) string {
	minSize := [...]int{0, 2, 4, 8}
	maxSize := [...]int{100, 3, 7, 100}
	if c.Empty() {
		return ""
	}
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
