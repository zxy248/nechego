package command

import (
	"nechego/chat"
	"nechego/game"
	"nechego/services"
	"regexp"
	"strings"
)

type RemoveCommandHandler struct {
	Universe *game.Universe
}

var removeCommandRe = regexp.MustCompile("^!(удалить|убрать) ([^\\|]+)")

func (h *RemoveCommandHandler) Match(m *chat.Message) services.Request {
	match := removeCommandRe.FindStringSubmatch(m.Text)
	if match == nil {
		return nil
	}
	def := strings.ToLower(strings.TrimSpace(match[2]))
	return &RemoveCommand{m, def, services.GetWorld(h.Universe, m.Group.ID)}
}

type RemoveCommand struct {
	Message    *chat.Message
	Definition string
	World      *game.World
}

func (r *RemoveCommand) Process() error {
	services.Do(r.World, func(w *game.World) {
		w.Commands.Remove(r.Definition)
	})
	return r.Message.Reply("❌ Команда удалена.")
}
