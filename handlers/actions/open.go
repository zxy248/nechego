package actions

import (
	"fmt"
	"html"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	tu "nechego/teleutil"
	"nechego/token"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

type Open struct {
	Universe *game.Universe
}

var openRe = handlers.Regexp("^!(прочитать|прочесть|открыть|распаковать) ([0-9]+)")

func (h *Open) Match(c tele.Context) bool {
	return openRe.MatchString(c.Text())
}

func (h *Open) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	k, ok := openKey(c.Text())
	if !ok {
		return c.Send(format.ChooseBox)
	}
	i, ok := user.Inventory.ByKey(k)
	if !ok {
		return c.Send(format.ItemNotFound)
	}
	switch x := i.Value.(type) {
	case *token.Letter:
		return h.handleLetter(c, x)
	case *item.Box:
		user.Inventory.Remove(i)
		return h.handleBox(c, user, x)
	default:
		return c.Send(format.ChooseBox)
	}
}

func (h *Open) handleLetter(c tele.Context, l *token.Letter) error {
	s := formatLetter(l)
	return c.Send(s, tele.ModeHTML)
}

func (h *Open) handleBox(c tele.Context, u *game.User, b *item.Box) error {
	u.Inventory.Add(b.Content)
	m := tu.Link(c, u)
	s := formatOpen(m, b)
	return c.Send(s, tele.ModeHTML)
}

func formatLetter(l *token.Letter) string {
	s := fmt.Sprintf("<b>✉️ Письмо</b> <i>(автор: <b>%s</b>)</i>\n", l.Author)
	s += fmt.Sprintf("<blockquote>%s</blockquote>", html.EscapeString(l.Text))
	return s
}

func formatOpen(who string, b *item.Box) string {
	return fmt.Sprintf("📦 %s открывает коробку. Внутри оказывается %s.",
		format.Name(who), format.Item(b.Content))
}

func openKey(s string) (k int, ok bool) {
	m := openRe.FindStringSubmatch(s)[2]
	k, err := strconv.Atoi(m)
	return k, err == nil
}
