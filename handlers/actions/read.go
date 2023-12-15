package actions

import (
	"fmt"
	"html"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"
	"nechego/token"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

type Read struct {
	Universe *game.Universe
}

var readRe = handlers.Regexp("^!прочитать ([0-9]+)")

func (h *Read) Match(c tele.Context) bool {
	return readRe.MatchString(c.Text())
}

func (h *Read) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	k, ok := letterKey(c.Text())
	if !ok {
		return c.Send(format.ChooseLetter)
	}
	i, ok := user.Inventory.ByKey(k)
	if !ok {
		return c.Send(format.ItemNotFound)
	}
	l, ok := i.Value.(*token.Letter)
	if !ok {
		return c.Send(format.ChooseLetter)
	}
	s := formatLetter(l)
	return c.Send(s, tele.ModeHTML)
}

func formatLetter(l *token.Letter) string {
	s := fmt.Sprintf("<b>✉️ Письмо</b> <i>(автор: <b>%s</b>)</i>\n", l.Author)
	s += fmt.Sprintf("<blockquote>%s</blockquote>", html.EscapeString(l.Text))
	return s
}

func letterKey(s string) (k int, ok bool) {
	m := readRe.FindStringSubmatch(s)[1]
	k, err := strconv.Atoi(m)
	return k, err == nil
}
