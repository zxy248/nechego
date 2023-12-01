package market

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Job struct {
	Universe *game.Universe
}

var jobRe = handlers.Regexp("^!(рохля|работа)")

func (h *Job) Match(c tele.Context) bool {
	return jobRe.MatchString(c.Text())
}

func (h *Job) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	const hr = 2
	const d = hr * time.Hour
	if !world.Market.Shift.Start(user.ID, d) {
		return c.Send(format.CannotGetJob)
	}
	l := tu.Link(c, user)
	s := format.GetJob(l, hr)
	return c.Send(s, tele.ModeHTML)
}
