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

func (h *Reputation) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!реп")
}

func (h *Reputation) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	who := tu.Link(c, user)
	r := format.Reputation{
		Score:  user.Reputation.Score(),
		Factor: user.ReputationFactor,
	}
	s := r.String(who)
	return c.Send(s, tele.ModeHTML)
}

type UpdateReputation struct {
	Universe *game.Universe
}

func (h *UpdateReputation) Match(c tele.Context) bool {
	_, ok := reputationDirection(c.Text())
	return ok
}

func (h *UpdateReputation) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	target, ok := tu.RepliedUser(c, world)
	if !ok || !canUpdateReputation(user, target) {
		return nil
	}

	d, _ := reputationDirection(c.Text())
	target.Reputation = target.Reputation.Update(user.ID, d)

	who := tu.Link(c, target)
	r := format.Reputation{
		Score:  target.Reputation.Score(),
		Factor: target.ReputationFactor,
	}
	s := r.Updated(who, d)
	return c.Send(s, tele.ModeHTML)
}

func reputationDirection(s string) (d reputation.Direction, ok bool) {
	switch s {
	case "+":
		return reputation.Up, true
	case "-":
		return reputation.Down, true
	}
	return 0, false
}

func canUpdateReputation(u1, u2 *game.User) bool {
	return time.Since(u2.Reputation.Last(u1.ID).Time) > 4*time.Hour
}
