package command

import (
	"nechego/chat"
	"nechego/commands"
	"nechego/game"
	"nechego/services"
	"regexp"
	"strings"
)

type AddCommandHandler struct {
	Universe *game.Universe
}

var addRe = regexp.MustCompile("^!добавить ([^\\|]+)\\|?(.*)")

func (h *AddCommandHandler) Match(m *chat.Message) services.Request {
	match := addRe.FindStringSubmatch(m.Text)
	if match == nil {
		return nil
	}
	def := strings.ToLower(strings.TrimSpace(match[1]))
	sub := strings.TrimSpace(match[2])
	return &AddCommand{m, def, sub, services.GetWorld(h.Universe, m.Group.ID)}
}

type AddCommand struct {
	Message      *chat.Message
	Definition   string
	Substitution string
	World        *game.World
}

func (r *AddCommand) Process() error {
	services.Do(r.World, func(w *game.World) {
		if w.Commands == nil {
			w.Commands = commands.Commands{}
		}
		cmd := commands.Command{
			Message: r.Substitution,
			Photo:   r.Message.File,
		}
		w.Commands.Add(r.Definition, cmd)
	})
	return r.Message.Reply("✅ Команда добавлена.")
}
