package command

import (
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Remove struct {
	Universe *game.Universe
}

var removeRe = handlers.NewRegexp(removePattern)

func (h *Remove) Match(c tele.Context) bool {
	return removeRe.MatchString(c.Text())
}

func (h *Remove) Handle(c tele.Context) error {
	world := tu.Lock(c, h.Universe)
	defer world.Unlock()

	m := removeRe.FindStringSubmatch(c.Text())
	d := sanitizeDefinition(m[2])
	world.Commands.Remove(d)
	return c.Send("❌ Команда удалена.")
}
