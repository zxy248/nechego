package command

import (
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type Remove struct {
	Universe *game.Universe
}

func (h *Remove) Match(s string) bool {
	_, ok := parseRemove(s)
	return ok
}

func (h *Remove) Handle(c tele.Context) error {
	def, _ := parseRemove(c.Text())
	tu.ContextWorld(c, h.Universe, func(w *game.World) {
		ensureCommands(w)
		w.Commands.Remove(def)
	})
	return c.Send("❌ Команда удалена.")
}

var removeRe = handlers.Regexp("^!(удалить|убрать) ([^\\|]+)")

func parseRemove(s string) (def string, ok bool) {
	m := removeRe.FindStringSubmatch(s)
	if m == nil {
		return "", false
	}
	return strings.ToLower(strings.TrimSpace(m[2])), true
}
