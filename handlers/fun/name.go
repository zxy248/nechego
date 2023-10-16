package fun

import (
	"html"
	"nechego/format"
	"nechego/handlers"
	"nechego/handlers/parse"
	tu "nechego/teleutil"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"
)

type Name struct{}

func (h *Name) Match(s string) bool {
	_, ok := parseName(s)
	return ok
}

func (h *Name) Handle(c tele.Context) error {
	u, ok := tu.Reply(c)
	if !ok {
		u = c.Sender()
	}
	name, _ := parseName(c.Message().Text)
	if longName(name) {
		return c.Send(format.LongName(maxNameLength), tele.ModeHTML)
	}
	if err := authorizeName(c, u); err != nil {
		return err
	}
	if err := setName(c, u, name); err != nil {
		return c.Send(format.CannotSetName, tele.ModeHTML)
	}
	return c.Send(format.NameSet(name), tele.ModeHTML)
}

func parseName(s string) (name string, ok bool) {
	if !parse.Seq(
		parse.Match("!имя"),
		parse.Str(parse.Assign(&name)),
	)(s) {
		return "", false
	}
	return html.EscapeString(name), true
}

const maxNameLength = 16

func longName(n string) bool {
	return utf8.RuneCountInString(n) > maxNameLength
}

func authorizeName(c tele.Context, u *tele.User) error {
	return tu.Promote(c, tu.Member(c, u))
}

func setName(c tele.Context, u *tele.User, n string) error {
	return c.Bot().SetAdminTitle(c.Chat(), u, n)
}

type CheckName struct{}

func (h *CheckName) Match(s string) bool {
	return handlers.MatchPrefix("!имя", s)
}

func (h *CheckName) Handle(c tele.Context) error {
	name := tu.Mention(c, c.Sender())
	return c.Send(format.YourName(name), tele.ModeHTML)
}
