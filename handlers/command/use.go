package command

import (
	"context"
	"math/rand/v2"
	"strings"

	"github.com/zxy248/nechego/data"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Use struct {
	Queries *data.Queries
}

func (h *Use) Match(c tele.Context) bool {
	ctx := context.Background()
	cmds, err := h.Queries.ListCommands(ctx, c.Chat().ID)
	if err != nil {
		return false
	}
	def := sanitizeDefinition(c.Text())
	_, ok := matchCommand(cmds, def)
	return ok
}

func (h *Use) Handle(c tele.Context) error {
	ctx := context.Background()
	cmds, err := h.Queries.ListCommands(ctx, c.Chat().ID)
	if err != nil {
		return err
	}
	def := sanitizeDefinition(c.Text())
	cmd, _ := matchCommand(cmds, def)
	if cmd.SubstitutionPhoto != "" {
		f := tele.File{FileID: cmd.SubstitutionPhoto}
		p := &tele.Photo{File: f, Caption: cmd.SubstitutionText}
		return c.Send(p)
	}
	return c.Send(cmd.SubstitutionText)
}

func matchCommand(cs []data.Command, s string) (c data.Command, ok bool) {
	var n int
	var r []data.Command
	for _, c := range cs {
		if len(c.Definition) > n && strings.Contains(s, c.Definition) {
			r = append(r, c)
			n = len(c.Definition)
		}
	}
	if r != nil {
		return r[rand.N(len(r))], true
	}
	return data.Command{}, false
}
