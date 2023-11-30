package fun

import (
	"html"
	"nechego/format"
	"nechego/handlers"
	tu "nechego/teleutil"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"
)

type Name struct{}

var nameRe = handlers.Regexp("^!имя (.+)")

func (h *Name) Match(c tele.Context) bool {
	return nameRe.MatchString(c.Text())
}

func (h *Name) Handle(c tele.Context) error {
	u, ok := tu.Reply(c)
	if !ok {
		u = c.Sender()
	}
	name := parseName(c.Text())
	if !validNameLength(name) {
		return c.Send(format.LongName(maxNameLength))
	}
	if err := promoteUser(c, u); err != nil {
		return err
	}
	if err := setName(c, u, name); err != nil {
		return c.Send(format.CannotSetName)
	}
	return c.Send(format.NameSet(name), tele.ModeHTML)
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

type CheckName struct{}

func (h *CheckName) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!имя")
}

func (h *CheckName) Handle(c tele.Context) error {
	l := tu.Link(c, c.Sender())
	return c.Send(format.YourName(l), tele.ModeHTML)
}
