package fun

import (
	"nechego/format"
	"nechego/game"
	"nechego/game/reputation"
	"nechego/handlers"
	tu "nechego/teleutil"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Reputation struct {
	Universe *game.Universe
}

func (h *Reputation) Match(s string) bool {
	return handlers.HasPrefix(s, "!реп")
}

func (h *Reputation) Handle(c tele.Context) error {
	return handlers.HandleWorld(c, h.Universe, h)
}

func (_ *Reputation) HandleWorld(c tele.Context, w *game.World) error {
	u := tu.CurrentUser(c, w)
	who := tu.Link(c, u)
	score := u.Reputation.Score()
	low := w.MinReputation()
	high := w.MaxReputation()
	r := format.ReputationSuffix(score, low, high)
	return c.Send(format.ReputationScore(who, r), tele.ModeHTML)
}

type UpdateReputation struct {
	Universe *game.Universe
}

func (h *UpdateReputation) Match(s string) bool {
	_, ok := reputationDirection(s)
	return ok
}

func (h *UpdateReputation) Handle(c tele.Context) error {
	return handlers.HandleWorld(c, h.Universe, h)
}

func (_ *UpdateReputation) HandleWorld(c tele.Context, w *game.World) error {
	u := tu.CurrentUser(c, w)
	v, ok := tu.RepliedUser(c, w)
	if !ok || !canUpdateReputation(u, v) {
		return nil
	}

	d, _ := reputationDirection(c.Message().Text)
	v.Reputation = v.Reputation.Update(u.TUID, d)

	who := tu.Link(c, v)
	score := v.Reputation.Score()
	low := w.MinReputation()
	high := w.MaxReputation()
	r := format.ReputationSuffix(score, low, high)
	return c.Send(format.ReputationUpdated(who, r, d), tele.ModeHTML)
}

func reputationDirection(s string) (d reputation.Direction, ok bool) {
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
