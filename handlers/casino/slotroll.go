package casino

import (
	"nechego/format"
	"nechego/game"
	"nechego/slot"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type SlotRoll struct {
	Universe *game.Universe
}

func (h *SlotRoll) Match(c tele.Context) bool {
	d := c.Message().Dice
	return d != nil && d.Type == "🎰"
}

func (h *SlotRoll) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	bet := user.SlotBet
	if bet == 0 {
		return nil
	}
	if !user.Balance().Spend(bet) {
		return c.Send(format.NoMoney)
	}

	m := tu.Link(c, user)
	v := c.Message().Dice.Value
	if p := slot.Prize(v, bet); p > 0 {
		user.Balance().Add(p)
		s := format.SlotWin(m, p)
		return c.Send(s, tele.ModeHTML)
	}
	s := format.SlotRoll(m, bet)
	return c.Send(s, tele.ModeHTML)
}
