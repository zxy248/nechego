package fun

import (
	"nechego/handlers"
	"nechego/handlers/fun/clock"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Clock struct{}

var clockRe = handlers.Regexp("!время до ([0-9:]+)")

func (*Clock) Match(c tele.Context) bool {
	_, ok := parseClock(c.Text())
	return ok
}

func (*Clock) Handle(c tele.Context) error {
	to, _ := parseClock(c.Text())
	from := clock.FromTime(time.Now())
	return c.Send(to.Sub(from).String())
}

func parseClock(s string) (c clock.Clock, ok bool) {
	m := clockRe.FindStringSubmatch(s)
	if m == nil {
		return clock.Clock{}, false
	}
	c, err := clock.FromString(m[1])
	return c, err == nil
}
