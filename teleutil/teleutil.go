package teleutil

import (
	"fmt"
	"github.com/zxy248/nechego/game"
	"html"
	"strings"

	tele "gopkg.in/zxy248/telebot.v3"
)

func Name(m *tele.ChatMember) string {
	name := m.Title
	if name == "" {
		return strings.TrimSpace(m.User.FirstName + " " + m.User.LastName)
	}
	return name
}

func Link(c tele.Context, who any) string {
	var m *tele.ChatMember
	switch x := who.(type) {
	case *tele.ChatMember:
		m = x
	case tele.Recipient:
		m = Member(c, x)
	case int64:
		m = Member(c, tele.ChatID(x))
	default:
		panic(fmt.Sprintf("unexpected type %T", x))
	}
	const format = `<a href="tg://user?id=%d">%s</a>`
	return fmt.Sprintf(format, m.User.ID, html.EscapeString(Name(m)))
}

func Member(c tele.Context, r tele.Recipient) *tele.ChatMember {
	m, err := c.Bot().ChatMemberOf(c.Chat(), r)
	if err != nil {
		panic("cannot get chat member")
	}
	return m
}

func Promote(c tele.Context, m *tele.ChatMember) error {
	if Admin(m) {
		return nil
	}
	m.Rights.CanBeEdited = true
	m.Rights.CanManageChat = true
	return c.Bot().Promote(c.Chat(), m)
}

func Admin(m *tele.ChatMember) bool {
	return m.Role == tele.Administrator || m.Role == tele.Creator
}

func Left(m *tele.ChatMember) bool {
	return m.Role == tele.Kicked || m.Role == tele.Left
}

func Reply(c tele.Context) *tele.User {
	if c.Message().IsReply() && !c.Message().ReplyTo.Sender.IsBot {
		return c.Message().ReplyTo.Sender
	}
	return nil
}

func Lock(c tele.Context, u *game.Universe) *game.World {
	world, err := u.World(c.Chat().ID)
	if err != nil {
		panic(fmt.Sprintf("cannot get world: %s", err))
	}
	world.Lock()
	return world
}

func MessageForwarded(m *tele.Message) bool {
	return m.OriginalUnixtime != 0
}

func SuperGroup(c tele.Context) bool {
	return c.Chat().Type == tele.ChatSuperGroup
}
