package profile

import (
	"nechego/avatar"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Profile struct {
	Universe *game.Universe
	Avatars  *avatar.Storage
}

var profileRe = handlers.Regexp("^!(профиль|стат)")

func (h *Profile) Match(c tele.Context) bool {
	return profileRe.MatchString(c.Text())
}

func (h *Profile) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if u, ok := tu.Reply(c); ok {
		user = world.User(u.ID)
	}
	out := format.Profile(user)
	if a, ok := h.Avatars.Get(user.ID); ok {
		a.Caption = out
		return c.Send(a, tele.ModeHTML)
	}
	return c.Send(out, tele.ModeHTML)
}
