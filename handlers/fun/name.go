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
	name, _ := parseName(c.Message().Text)
	if nameTooLong(name) {
		return c.Send(format.NameTooLong(maxNameLength), tele.ModeHTML)
	}
	if err := giveNameRights(c); err != nil {
		return err
	}
	if err := setName(c, name); err != nil {
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

func nameTooLong(n string) bool {
	return utf8.RuneCountInString(n) > maxNameLength
}

func giveNameRights(c tele.Context) error {
	return tu.Promote(c, tu.Member(c, c.Sender()))
}

func setName(c tele.Context, n string) error {
	return c.Bot().SetAdminTitle(c.Chat(), c.Sender(), n)
}

type CheckName struct{}

func (h *CheckName) Match(s string) bool {
	return handlers.MatchPrefix("!имя", s)
}

func (h *CheckName) Handle(c tele.Context) error {
	name := tu.Mention(c, c.Sender())
	return c.Send(format.YourName(name), tele.ModeHTML)
}
