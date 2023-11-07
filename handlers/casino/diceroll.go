package casino

import (
	"nechego/format"
	"nechego/game"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type DiceRoll struct {
	Universe *game.Universe
}

func (h *DiceRoll) Match(c tele.Context) bool {
	d := c.Message().Dice
	return d != nil && d.Type == "🎲"
}

func (h *DiceRoll) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	game := world.Casino.Game()
	if !game.Going() || !game.Verify(user.ID) {
		return nil
	}

	v := c.Message().Dice.Value
	r := game.Finish(v)
	user.Balance().Add(r.Prize)
	s := format.DiceGameResult(r)
	return c.Send(s, tele.ModeHTML)
}
