package pictures

import (
	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Casper struct {
	Path string
}

var casperRe = handlers.NewRegexp("^!касп[ие]р")

func (h *Casper) Match(c tele.Context) bool {
	return casperRe.MatchString(c.Text())
}

func (h *Casper) Handle(c tele.Context) error {
	f, err := randomFile(h.Path)
	if err != nil {
		return err
	}
	return c.Send(&tele.Photo{File: tele.FromDisk(f)})
}
