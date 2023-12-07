package casino

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

type Slot struct {
	Universe *game.Universe
	MinBet   int
}

var slotRe = handlers.Regexp("^!(слот|ставка|казино) ([0-9]+)")

func (h *Slot) Match(c tele.Context) bool {
	return slotRe.MatchString(c.Text())
}

func (h *Slot) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	bet := slotBet(c.Text())
	if bet < h.MinBet {
		s := format.MinBet(h.MinBet)
		return c.Send(s, tele.ModeHTML)
	}
	user.SlotBet = bet

	m := tu.Link(c, user)
	s := format.BetSet(m, bet)
	return c.Send(s, tele.ModeHTML)
}

func slotBet(s string) int {
	m := slotRe.FindStringSubmatch(s)[2]
	n, _ := strconv.Atoi(m)
	return n
}
