package teleutil

import (
	"nechego/format"
	"regexp"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func Name(m *tele.ChatMember) string {
	name := m.Title
	if name == "" {
		return strings.TrimSpace(m.User.FirstName + " " + m.User.LastName)
	}
	return name
}

func Mention(c tele.Context, gid, uid int64) string {
	member, err := c.Bot().ChatMemberOf(tele.ChatID(gid), tele.ChatID(uid))
	if err != nil {
		return "❔"
	}
	return format.Mention(uid, Name(member))
}

func Args(c tele.Context, re *regexp.Regexp) []string {
	return re.FindStringSubmatch(c.Message().Text)
}
