package command

import (
	"nechego/commands"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Add struct {
	Universe *game.Universe
}

var addRe = handlers.NewRegexp(addPattern)

func (h *Add) Match(c tele.Context) bool {
	return addRe.MatchString(c.Text())
}

func (h *Add) Handle(c tele.Context) error {
	world := tu.Lock(c, h.Universe)
	defer world.Unlock()

	m := addRe.FindStringSubmatch(c.Text())
	d := sanitizeDefinition(m[1])
	s := sanitizeSubstitution(m[2])
	x := commands.Command{
		Message: s,
		Photo:   photoFileID(c),
	}
	world.Commands.Add(d, x)
	return c.Send("✅ Команда добавлена.")
}

func photoFileID(c tele.Context) string {
	var s string
	if p := c.Message().Photo; p != nil {
		s = p.FileID
	}
	return s
}
