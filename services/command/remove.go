package command

import (
	"nechego/chat"
	"nechego/game"
	"nechego/services"
)

type Remove struct {
	Universe *game.Universe
}

var removeRe = services.Regexp(removePattern)

func (h *Remove) Match(m *chat.Message) services.Request {
	match := removeRe.FindStringSubmatch(m.Text)
	if match == nil {
		return nil
	}
	def := sanitizeDefinition(match[2])
	s := services.GroupState(h.Universe, m)
	return &remove{s, m, def}
}

type remove struct {
	state      *services.State
	message    *chat.Message
	definition string
}

func (r *remove) Process() error {
	r.state.Do(func(w *game.World) {
		w.Commands.Remove(r.definition)
	})
	return r.message.Reply("❌ Команда удалена.")
}
