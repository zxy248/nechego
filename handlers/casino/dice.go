package casino

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

type Dice struct {
	Universe *game.Universe
	MinBet   int
}

var diceRe = handlers.Regexp("^!кости ([0-9]+)")

func (h *Dice) Match(c tele.Context) bool {
	return diceRe.MatchString(c.Text())
}

func (h *Dice) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	casino := world.Casino
	bet := diceBet(c.Text())
	if bet < h.MinBet {
		s := format.MinBet(h.MinBet)
		return c.Send(s, tele.ModeHTML)
	}
	if casino.Game().Going() {
		return c.Send(format.GameGoing)
	}
	if !user.Balance().Spend(bet) {
		return c.Send(format.NoMoney)
	}

	roll := rollDiceFunc(c)
	timeout := timeoutFunc(c, bet)
	if err := casino.PlayDice(user.TUID, bet, roll, timeout); err != nil {
		return err
	}

	m := tu.Link(c, user)
	s := format.DiceGame(m, bet, casino.Timeout)
	return c.Send(s, tele.ModeHTML)
}

func diceBet(s string) int {
	m := diceRe.FindStringSubmatch(s)[1]
	n, _ := strconv.Atoi(m)
	return n
}

func rollDiceFunc(c tele.Context) game.RollDiceFunc {
	return func() (score int, err error) {
		opt := &tele.SendOptions{ReplyTo: c.Message()}
		m, err := tele.Cube.Send(c.Bot(), c.Chat(), opt)
		if err != nil {
			return 0, err
		}
		return m.Dice.Value, nil
	}
}

func timeoutFunc(c tele.Context, bet int) func() {
	return func() {
		s := format.DiceTimeout(bet)
		c.Send(s, tele.ModeHTML)
	}
}
