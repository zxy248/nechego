package farm

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"
	"nechego/valid"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type Name struct {
	Universe *game.Universe
}

var nameRe = handlers.Regexp("!назвать ферму (.+)")

func (h *Name) Match(c tele.Context) bool {
	return nameRe.MatchString(c.Text())
}

func (h *Name) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	name := nameRe.FindStringSubmatch(c.Text())[1]
	if !valid.Name(name) {
		return c.Send(format.BadFarmName)
	}
	user.Farm.Name = strings.Title(name)
	s := format.FarmNamed(tu.Mention(c, user), user.Farm)
	return c.Send(s, tele.ModeHTML)
}
