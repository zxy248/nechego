package command

import (
	"nechego/chat"
	"nechego/commands"
	"nechego/game"
	"nechego/services"
	"strings"
)

type UseCommandHandler struct {
	Universe *game.Universe
}

func (h *UseCommandHandler) Match(m *chat.Message) services.Request {
	w := services.GetWorld(h.Universe, m.Group.ID)
	var cmd commands.Command
	var ok bool
	services.Do(w, func(w *game.World) {
		cmd, ok = w.Commands.Match(strings.ToLower(m.Text))
	})
	if !ok {
		return nil
	}
	return &UseCommand{m, w, cmd}
}

type UseCommand struct {
	Message *chat.Message
	World   *game.World
	Command commands.Command
}

func (r *UseCommand) Process() error {
	if r.Command.HasPhoto() {
		p := chat.Photo{File: r.Command.Photo}
		return r.Message.ReplyPhoto(r.Command.Message, p)
	}
	return r.Message.Reply(r.Command.Message)
}
