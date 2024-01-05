package command

import (
	"nechego/game"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Use struct {
	Universe *game.Universe
}

func (h *Use) Match(c tele.Context) bool {
	world := tu.Lock(c, h.Universe)
	defer world.Unlock()

	d := sanitizeDefinition(c.Text())
	_, ok := world.Commands.Match(d)
	return ok
}

func (h *Use) Handle(c tele.Context) error {
	world := tu.Lock(c, h.Universe)
	defer world.Unlock()

	d := sanitizeDefinition(c.Text())
	cmd, _ := world.Commands.Match(d)
	if cmd.HasPhoto() {
		f := tele.File{FileID: cmd.Photo}
		p := &tele.Photo{File: f, Caption: cmd.Message}
		return c.Send(p)
	}
	return c.Send(cmd.Message)
}
