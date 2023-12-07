package actions

import (
	"nechego/fishing"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Records struct {
	Universe *game.Universe
}

var recordsRe = handlers.Regexp("^!рекорды")

func (h *Records) Match(c tele.Context) bool {
	return recordsRe.MatchString(c.Text())
}

func (h *Records) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if len(world.History.Entries) == 0 {
		return c.Send(format.NoFishingRecords)
	}
	price := world.History.Top(fishing.Price, 10)
	weight := world.History.Top(fishing.Weight, 3)
	length := world.History.Top(fishing.Length, 3)
	s := format.FishingRecords(price, weight, length)
	return c.Send(s, tele.ModeHTML)
}
