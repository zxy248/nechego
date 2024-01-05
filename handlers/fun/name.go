package fun

import (
	"fmt"
	"html"
	"nechego/handlers"
	tu "nechego/teleutil"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"
)

type Name struct{}

var nameRe = handlers.NewRegexp("^!имя (.+)")

func (h *Name) Match(c tele.Context) bool {
	return nameRe.MatchString(c.Text())
}

func (h *Name) Handle(c tele.Context) error {
	u := tu.Reply(c)
	if u == nil {
		u = c.Sender()
	}
	name := parseName(c.Text())
	if !validNameLength(name) {
		return c.Send(nameLengthExceeded(maxNameLength))
	}
	if err := promoteUser(c, u); err != nil {
		return err
	}
	if err := setName(c, u, name); err != nil {
		return c.Send(setNameFail())
	}
	return c.Send(setNameSuccess(name), tele.ModeHTML)
}

func parseName(s string) string {
	return html.EscapeString(nameRe.FindStringSubmatch(s)[1])
}

const maxNameLength = 16

func validNameLength(s string) bool {
	return utf8.RuneCountInString(s) <= maxNameLength
}

func promoteUser(c tele.Context, u *tele.User) error {
	return tu.Promote(c, tu.Member(c, u))
}

func setName(c tele.Context, u *tele.User, name string) error {
	return c.Bot().SetAdminTitle(c.Chat(), u, name)
}

func setNameSuccess(name string) string {
	return fmt.Sprintf("Имя <b>%s</b> установлено ✅", name)
}

func setNameFail() string {
	return "⚠️ Не удалось установить имя."
}

func nameLengthExceeded(max int) string {
	return fmt.Sprintf("⚠️ Максимальная длина имени %d символов.", max)
}

type CheckName struct{}

func (h *CheckName) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!имя")
}

func (h *CheckName) Handle(c tele.Context) error {
	l := tu.Link(c, c.Sender())
	s := fmt.Sprintf("Ваше имя: <b>%s</b> 🔖", l)
	return c.Send(s, tele.ModeHTML)
}
