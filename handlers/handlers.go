package handlers

import (
	"nechego/game"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type WorldHandler interface {
	HandleWorld(c tele.Context, w *game.World) error
}

func HandleWorld(c tele.Context, u tu.WorldGetter, h WorldHandler) error {
	var err error
	tu.ContextWorld(c, u, func(w *game.World) {
		err = h.HandleWorld(c, w)
	})
	return err
}
