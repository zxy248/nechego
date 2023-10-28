package command

import (
	"nechego/commands"
	"nechego/game"
	tu "nechego/teleutil"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type Use struct {
	Universe *game.Universe
}

func (h *Use) Match(c tele.Context) bool {
	_, ok := h.parseUse(c)
	return ok
}

func (h *Use) Handle(c tele.Context) error {
	cmd, _ := h.parseUse(c)
	if cmd.HasPhoto() {
		f := tele.File{FileID: cmd.Photo}
		p := tele.Photo{File: f, Caption: cmd.Message}
		return c.Send(&p)
	}
	return c.Send(cmd.Message)
}

func (h *Use) parseUse(c tele.Context) (cmd commands.Command, ok bool) {
	tu.ContextWorld(c, h.Universe, func(w *game.World) {
		cmd, ok = w.Commands.Match(strings.ToLower(c.Text()))
	})
	return
}
