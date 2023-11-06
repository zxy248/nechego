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

	n := farmName(c.Text())
	if n == "" {
		return c.Send(format.BadFarmName)
	}
	user.Farm.Name = n
	s := format.FarmNamed(tu.LinkSender(c), n)
	return c.Send(s, tele.ModeHTML)
}

func farmName(s string) string {
	n := nameRe.FindStringSubmatch(s)[1]
	if !valid.Name(n) {
		return ""
	}
	return strings.Title(n)
}
