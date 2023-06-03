package fun

import (
	"nechego/handlers/fun/clock"
	"nechego/handlers/parse"

	tele "gopkg.in/telebot.v3"
)

type Time struct{}

func (*Time) Match(s string) bool {
	_, ok := timeCommand(s)
	return ok
}

func (*Time) Handle(c tele.Context) error {
	c1, _ := timeCommand(c.Message().Text)
	c2 := clock.Now()
	d := c1.Sub(c2)
	return c.Send(d.String())
}

func timeCommand(s string) (c clock.Clock, ok bool) {
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
