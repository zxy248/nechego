package casino

import (
	"nechego/format"
	"nechego/game"
	"nechego/slot"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Roll struct {
	Universe *game.Universe
}

func (h *Roll) Match(c tele.Context) bool {
	return c.Message().Dice != nil
}

func (h *Roll) Handle(c tele.Context) error {
	switch c.Message().Dice.Type {
	case "🎲":
		return h.handleDiceRoll(c)
	case "🎰":
		return h.handleSlotRoll(c)
	}
	return nil
}

func (h *Roll) handleDiceRoll(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	game := world.Casino.Game()
	if !game.Going() || !game.Verify(user.TUID) {
		return nil
	}

	res := game.Finish(c.Message().Dice.Value)
	user.Balance().Add(res.Prize)
	return c.Send(format.DiceGameResult(res), tele.ModeHTML)
}

func (h *Roll) handleSlotRoll(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	bet := user.SlotBet
	if bet == 0 {
		return nil
	}
	if !user.Balance().Spend(bet) {
		return c.Send(format.NoMoney)
	}

	mention := tu.Mention(c, user)
	value := c.Message().Dice.Value
	if prize, ok := slot.Prize(value, bet); ok {
		user.Balance().Add(prize)
		return c.Send(format.SlotWin(mention, prize), tele.ModeHTML)
	}
	return c.Send(format.SlotRoll(mention, bet), tele.ModeHTML)
}
