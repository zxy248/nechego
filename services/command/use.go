package command

import (
	"nechego/chat"
	"nechego/commands"
	"nechego/game"
	"nechego/services"
)

type Use struct {
	Universe *game.Universe
}

func (h *Use) Match(m *chat.Message) services.Request {
	s := services.GroupState(h.Universe, m)
	var cmd commands.Command
	var ok bool
	s.Do(func(w *game.World) {
		cmd, ok = w.Commands.Match(sanitizeDefinition(m.Text))
	})
	if !ok {
		return nil
	}
	return &use{s, m, cmd}
}

type use struct {
	state   *services.State
	message *chat.Message
	command commands.Command
}

func (r *use) Process() error {
	if r.command.HasPhoto() {
		p := chat.Photo{File: r.command.Photo}
		return r.message.ReplyPhoto(r.command.Message, p)
	}
	return r.message.Reply(r.command.Message)
}
