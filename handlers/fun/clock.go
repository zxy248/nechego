package fun

import (
	"nechego/handlers/fun/clock"
	"nechego/handlers/parse"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Clock struct{}

func (*Clock) Match(s string) bool {
	_, ok := parseClock(s)
	return ok
}

func (*Clock) Handle(c tele.Context) error {
	given, _ := parseClock(c.Message().Text)
	now := clock.FromTime(time.Now())
	return c.Send(given.Sub(now).String())
}

func parseClock(s string) (c clock.Clock, ok bool) {
	var x string
	if !parse.Seq(
		parse.Match("!время"),
		parse.Match("до"),
		parse.Str(parse.Assign(&x)),
	)(s) {
		return clock.Clock{}, false
	}
	c, err := clock.FromString(x)
	return c, err == nil
}
