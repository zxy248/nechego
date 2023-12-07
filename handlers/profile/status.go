package profile

import (
	"fmt"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"
)

type Status struct {
	Universe  *game.Universe
	MaxLength int
}

var statusRe = handlers.Regexp("^!статус (.*)")

func (h *Status) Match(c tele.Context) bool {
	return statusRe.MatchString(c.Text())
}

func (h *Status) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	s := statusText(c.Text())
	if utf8.RuneCountInString(s) > h.MaxLength {
		const f = "💬 Максимальная длина статуса %d символов."
		return c.Send(fmt.Sprintf(f, h.MaxLength))
	}
	user.Status = s
	return c.Send("✅ Статус установлен.")
}

func statusText(s string) string {
	m := statusRe.FindStringSubmatch(s)
	return m[1]
}
