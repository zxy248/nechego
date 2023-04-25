package casino

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers/parse"
	tu "nechego/teleutil"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Dice struct {
	Universe *game.Universe
}

func (h *Dice) Match(s string) bool {
	_, ok := diceCommand(s)
	return ok
}

func (h *Dice) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	bet, ok := diceCommand(c.Text())
	if !ok {
		panic("bad dice command")
	}
	if min := 100; bet < min {
		return c.Send(format.MinBet(min), tele.ModeHTML)
	}
	if !user.Dice() {
		return c.Send(format.NoDice)
	}
	if world.Casino.Game().Going() {
		return c.Send(format.GameGoing)
	}
	if !user.Balance().Spend(bet) {
		return c.Send(format.NoMoney)
	}
	throw := diceThrowFunc(c)
	timeout := diceTimeoutFunc(c, bet)
	if err := world.Casino.PlayDice(user.TUID, bet, throw, timeout); err != nil {
		return err
	}
	mention := tu.Mention(c, c.Sender())
	seconds := int(world.Casino.Timeout / time.Second)
	return c.Send(format.DiceGame(mention, bet, seconds), tele.ModeHTML)
}

func diceCommand(s string) (bet int, ok bool) {
	ok = parse.Seq(parse.Match("!кости"), parse.Int(parse.Assign(&bet)))(s)
	return
}

func diceThrowFunc(c tele.Context) game.DiceThrowFunc {
	return func() (score int, err error) {
		m, err := tele.Cube.Send(c.Bot(), c.Chat(), nil)
		if err != nil {
			return 0, err
		}
		return m.Dice.Value, nil
	}
}

func diceTimeoutFunc(c tele.Context, bet int) func() {
	return func() { c.Send(format.DiceTimeout(bet), tele.ModeHTML) }
}
