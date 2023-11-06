package top

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Rating struct {
	Universe *game.Universe
}

var ratingRe = handlers.Regexp("^!(рейтинг|ммр|эло)")

func (h *Rating) Match(c tele.Context) bool {
	return ratingRe.MatchString(c.Text())
}

func (h *Rating) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	us := trim(world.SortedUsers(game.ByElo))
	s := format.TopRating(whoFunc(c), us)
	return c.Send(s, tele.ModeHTML)
}
