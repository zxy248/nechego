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

// mentionName returns the mentionName of the user by the name.
func mentionName(uid int64, name string) string {
	return fmt.Sprintf("[%s](tg://user?id=%d)", name, uid)
}

// mentionUser returns the mention of the user by his name.
func (a *App) mentionUser(u model.User) (string, error) {
	m, err := a.chatMember(u)
	if err != nil {
		return "", err
	}
	name := markdownEscaper.Replace(chatMemberName(m))
	return mentionName(u.UID, name), nil
}

func (a *App) mustMentionUser(u model.User) string {
	name, err := a.mentionUser(u)
	if err != nil {
		a.SugarLog().Errorw("can't mention the user", "user", u)
		return "Имя не найдено"
	}
	return name
}

func respondPlain(c tele.Context, out string) error {
	return c.Send(out)
}

func respondMarkdown(c tele.Context, out string) error {
	return c.Send(out, tele.ModeMarkdownV2)
}

func respondHTML(c tele.Context, out string) error {
	return c.Send(out, tele.ModeHTML)
}

func internalError(c tele.Context, err error) error {
	respondPlain(c, makeError("Ошибка сервера"))
	return err
}

func userError(c tele.Context, out string) error {
	return respondPlain(c, makeError(out))
}

func userErrorMarkdown(c tele.Context, out string) error {
	return respondMarkdown(c, makeError(out))
}
