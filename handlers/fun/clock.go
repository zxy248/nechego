package fun

import (
	"github.com/zxy248/nechego/handlers"
	"github.com/zxy248/nechego/handlers/fun/clock"
	"time"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Clock struct{}

var clockRe = handlers.NewRegexp("^!время до ([0-9:]+)")

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
