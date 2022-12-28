package app

import (
	"fmt"
	"nechego/input"
	"nechego/model"
	"strings"

	tele "gopkg.in/telebot.v3"
)

// sendSmallProfilePhoto sends a small profile photo of the sender.
func sendSmallProfilePhoto(c tele.Context) error {
	u, err := c.Bot().ChatByID(c.Sender().ID)
	if err != nil {
		return err
	}
	tf, err := c.Bot().FileByID(u.Photo.SmallFileID)
	if err != nil {
		return err
	}
	f, err := c.Bot().File(&tf)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(f)})
}

// sendLargeProfilePhoto sends a large profile photo of the sender.
func sendLargeProfilePhoto(c tele.Context) error {
	p, err := c.Bot().ProfilePhotosOf(c.Sender())
	if err != nil {
		return err
	}
	if len(p) < 1 {
		return nil
	}
	return c.Send(&p[0])
}

// chatMemberAbsent returns true if the member is kicked or left.
func chatMemberAbsent(m *tele.ChatMember) bool {
	return m.Role == tele.Kicked || m.Role == tele.Left
}

func displayedName(m *tele.ChatMember) string {
	name := m.Title
	if name == "" {
		name = m.User.FirstName + " " + m.User.LastName
	}
	return strings.TrimSpace(name)
}

func (a *App) chatMember(u model.User) (*tele.ChatMember, error) {
	m, err := a.bot.ChatMemberOf(tele.ChatID(u.GID), tele.ChatID(u.UID))
	if err != nil {
		return nil, err
	}
	return m, nil
}

func formatMention(uid int64, name string) string {
	return fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, uid, input.Sanitize(name))
}

func (a *App) mention(u model.User) string {
	m, err := a.chatMember(u)
	if err != nil {
		return "❔"
	}
	return formatMention(u.UID, displayedName(m))
}
