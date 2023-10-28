package command

import (
	"nechego/commands"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type Add struct {
	Universe *game.Universe
}

func (h *Add) Match(s string) bool {
	_, _, ok := parseAdd(s)
	return ok
}

func (h *Add) Handle(c tele.Context) error {
	def, sub, _ := parseAdd(c.Text())
	tu.ContextWorld(c, h.Universe, func(w *game.World) {
		ensureCommands(w)
		cmd := commands.Command{Message: sub}
		w.Commands.Add(def, populateCommand(cmd, c))
	})
	return c.Send("✅ Команда добавлена.")
}

var addRe = handlers.Regexp("^!добавить ([^\\|]+)\\|?(.*)")

func parseAdd(s string) (def, sub string, ok bool) {
	m := addRe.FindStringSubmatch(s)
	if m == nil {
		return "", "", false
	}
	return strings.ToLower(strings.TrimSpace(m[1])), strings.TrimSpace(m[2]), true
}

func ensureCommands(w *game.World) {
	if w.Commands == nil {
		w.Commands = commands.Commands{}
	}
}

func populateCommand(cmd commands.Command, c tele.Context) commands.Command {
	p := c.Message().Photo
	if p != nil {
		cmd.Photo = p.FileID
	}
	return cmd
}
