package pictures

import (
	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Basili struct {
	Path string
}

var basiliRe = handlers.NewRegexp("^!(муся|марс|(кот|кош)[а-я]* василия)")

func (h *Basili) Match(c tele.Context) bool {
	return basiliRe.MatchString(c.Text())
}

func (h *Basili) Handle(c tele.Context) error {
	name, err := randomFile(h.Path)
	if err != nil {
		return err
	}
	return c.Send(sendableFromFile(name))
}
