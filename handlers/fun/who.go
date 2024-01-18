package fun

import (
	"html"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Who struct {
	Universe *game.Universe
}

var whoRe = handlers.NewRegexp("^!кто(.*)")

func (h *Who) Match(c tele.Context) bool {
	return whoRe.MatchString(c.Text())
}

func (h *Who) Handle(c tele.Context) error {
	world := tu.Lock(c, h.Universe)
	defer world.Unlock()

	w := parseWho(c.Text())
	l := tu.Link(c, world.RandomUser())
	s := l + w
	return c.Send(s, tele.ModeHTML)
}

func parseWho(s string) string {
	return html.EscapeString(whoRe.FindStringSubmatch(s)[1])
}
