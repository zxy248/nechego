package fun

import (
	"regexp"

	"github.com/zxy248/nechego/teleutil"
	tele "gopkg.in/zxy248/telebot.v3"
)

type Avatar struct{}

var avatarRe = regexp.MustCompile("^!ава")

func (h *Avatar) Match(c tele.Context) bool {
	return avatarRe.MatchString(c.Text())
}

func (h *Avatar) Handle(c tele.Context) error {
	user := teleutil.Reply(c)
	if user == nil {
		user = c.Sender()
	}

	photos, err := c.Bot().ProfilePhotosOf(user)
	if err != nil || len(photos) == 0 {
		return nil
	}
	return c.Send(&photos[0])
}
