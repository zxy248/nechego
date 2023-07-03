package handlers

import (
	"fmt"
	"nechego/game"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type WorldHandler interface {
	HandleWorld(c tele.Context, w *game.World) error
}

func HandleWorld(c tele.Context, u *game.Universe, h WorldHandler) error {
	w, err := u.World(c.Chat().ID)
	if err != nil {
		return fmt.Errorf("HandleWorld: cannot get world: %s", err)
	}
	w.Lock()
	defer w.Unlock()
	return h.HandleWorld(c, w)
}

func CurrentUser(c tele.Context, w *game.World) *game.User {
	return w.UserByID(c.Sender().ID)
}

func RepliedUser(c tele.Context, w *game.World) (u *game.User, ok bool) {
	r, ok := tu.Reply(c)
	if !ok {
		return nil, false
	}
	return w.UserByID(r.ID), true
}
