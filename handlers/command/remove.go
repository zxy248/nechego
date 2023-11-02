package command

import (
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Remove struct {
	Universe *game.Universe
}

var removeRe = handlers.Regexp(removePattern)

func (h *Remove) Match(c tele.Context) bool {
	return removeRe.MatchString(c.Text())
}

func (h *Remove) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	m := removeRe.FindStringSubmatch(c.Text())
	d := sanitizeDefinition(m[2])
	world.Commands.Remove(d)
	return c.Send("❌ Команда удалена.")
}
