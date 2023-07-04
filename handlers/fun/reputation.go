package fun

import (
	"nechego/format"
	"nechego/game"
	"nechego/game/reputation"
	hh "nechego/handlers"
	tu "nechego/teleutil"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Reputation struct {
	Universe *game.Universe
}

func (h *Reputation) Match(s string) bool {
	return hh.MatchPrefix("!репутация", s)
}

func (h *Reputation) Handle(c tele.Context) error {
	return hh.HandleWorld(c, h.Universe, h)
}

func (_ *Reputation) HandleWorld(c tele.Context, w *game.World) error {
	u := tu.CurrentUser(c, w)
	mention := tu.Mention(c, u)
	total := u.Reputation.Total()
	return c.Send(format.ReputationTotal(mention, total), tele.ModeHTML)
}

type UpdateReputation struct {
	Universe *game.Universe
}

func (h *UpdateReputation) Match(s string) bool {
	_, ok := reputationDirection(s)
	return ok
}

func (h *UpdateReputation) Handle(c tele.Context) error {
	return hh.HandleWorld(c, h.Universe, h)
}

func (_ *UpdateReputation) HandleWorld(c tele.Context, w *game.World) error {
	u := tu.CurrentUser(c, w)
	v, ok := tu.RepliedUser(c, w)
	if !ok || !canUpdateReputation(u, v) {
		return nil
	}

	d, _ := reputationDirection(c.Message().Text)
	v.Reputation = v.Reputation.Update(u.TUID, d)

	mention := tu.Mention(c, v)
	total := v.Reputation.Total()
	return c.Send(format.ReputationUpdated(mention, total, d), tele.ModeHTML)
}

func reputationDirection(s string) (d reputation.Dir, ok bool) {
	switch s {
	case "-":
		return reputation.Down, true
	case "+":
		return reputation.Up, true
	}
	return 0, false
}

func canUpdateReputation(u, v *game.User) bool {
	return u != v && time.Since(v.Reputation.Last(u.TUID).Time) > 4*time.Hour
}
