package casino

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers/parse"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Slot struct {
	Universe *game.Universe
	MinBet   int
}

func (h *Slot) Match(s string) bool {
	_, ok := slotCommand(s)
	return ok
}

func (h *Slot) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	bet, ok := slotCommand(c.Text())
	if !ok {
		panic("bad slot command")
	}
	if bet < h.MinBet {
		return c.Send(format.MinBet(h.MinBet), tele.ModeHTML)
	}
	user.SlotBet = bet
	return c.Send(format.BetSet(tu.Mention(c, user), bet), tele.ModeHTML)
}

func slotCommand(s string) (bet int, ok bool) {
	ok = parse.Seq(
		parse.Prefix("!слот", "!ставка", "!казино"),
		parse.Int(parse.Assign(&bet)),
	)(s)
	return
}
