package farm

import (
	"nechego/farm"
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

	n := farmName(c.Text())
	if !setFarmName(user.Farm, n) {
		return c.Send(format.BadFarmName)
	}
	s := format.FarmNamed(tu.LinkSender(c), n)
	return c.Send(s, tele.ModeHTML)
}

func farmName(s string) string {
	n := nameRe.FindStringSubmatch(s)[1]
	return strings.Title(n)
}

func setFarmName(f *farm.Farm, n string) bool {
	if !valid.Name(n) {
		return false
	}
	f.Name = n
	return true
}
