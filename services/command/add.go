package command

import (
	"nechego/chat"
	"nechego/commands"
	"nechego/game"
	"nechego/services"
)

type Add struct {
	Universe *game.Universe
}

var addRe = services.Regexp(addPattern)

func (h *Add) Match(m *chat.Message) services.Request {
	match := addRe.FindStringSubmatch(m.Text)
	if match == nil {
		return nil
	}
	def := sanitizeDefinition(match[1])
	sub := sanitizeSubstitution(match[2])
	s := services.GroupState(h.Universe, m)
	return &add{s, m, def, sub}
}

type add struct {
	state        *services.State
	message      *chat.Message
	definition   string
	substitution string
}

func (r *add) Process() error {
	r.state.Do(func(w *game.World) {
		if w.Commands == nil {
			w.Commands = commands.Commands{}
		}
		cmd := commands.Command{
			Message: r.substitution,
			Photo:   r.message.File,
		}
		w.Commands.Add(r.definition, cmd)
	})
	return r.message.Reply("✅ Команда добавлена.")
}
