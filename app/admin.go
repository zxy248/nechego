package app

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

// !очистка
func (a *App) handleCleanup(c tele.Context) error {
	users, err := a.model.ListUsers(getGroup(c))
	if err != nil {
		return err
	}
	var out string
	for _, u := range users {
		m, err := a.chatMember(u)
		if err != nil {
			return err
		}
		if !chatMemberPresent(m) {
			out += fmt.Sprintf("Пользователь _%s_ удален\\.",
				markdownEscaper.Replace(chatMemberName(m)))
			a.model.DeleteUser(u)
		}
	}
	return c.Send(out, tele.ModeMarkdownV2)
}
