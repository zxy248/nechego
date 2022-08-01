package app

import (
	"fmt"
	"nechego/model"
	"strings"

	tele "gopkg.in/telebot.v3"
)

// sendSmallProfilePhoto sends a small profile photo of the sender.
func sendSmallProfilePhoto(c tele.Context) error {
	user, err := c.Bot().ChatByID(c.Sender().ID)
	if err != nil {
		return err
	}
	file, err := c.Bot().FileByID(user.Photo.SmallFileID)
	if err != nil {
		return err
	}
	f, err := c.Bot().File(&file)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(f)})
}

// sendLargeProfilePhoto sends a large profile photo of the sender.
func sendLargeProfilePhoto(c tele.Context) error {
	ps, err := c.Bot().ProfilePhotosOf(c.Sender())
	if err != nil {
		return err
	}
	if len(ps) < 1 {
		return nil
	}
	return c.Send(&ps[0])
}

// chatMemberPresent returns true if the member is not kicked or left.
func chatMemberPresent(m *tele.ChatMember) bool {
	if m.Role == tele.Kicked || m.Role == tele.Left {
		return false
	}
	return true
}

// chatMemberName returns the member's displayed name.
func chatMemberName(m *tele.ChatMember) string {
	name := m.Title
	if name == "" {
		name = m.User.FirstName + " " + m.User.LastName
	}
	return strings.TrimSpace(name)
}

// chatMember gets a chat member by GID and UID.
func (a *App) chatMember(u model.User) (*tele.ChatMember, error) {
	member, err := a.bot.ChatMemberOf(tele.ChatID(u.GID), tele.ChatID(u.UID))
	if err != nil {
		return nil, err
	}
	return member, nil
}

// mentionUser returns the mention of the user by his name.
func (a *App) mentionUser(u model.User) (HTML, error) {
	m, err := a.chatMember(u)
	if err != nil {
		return "", err
	}
	return mention(u.UID, chatMemberName(m)), nil
}

func (a *App) mustMentionUser(u model.User) HTML {
	out, err := a.mentionUser(u)
	if err != nil {
		panic(fmt.Errorf("can't mention the user: %v", err))
	}
	return out
}
