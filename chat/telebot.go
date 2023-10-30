package chat

import tele "gopkg.in/telebot.v3"

func Convert(c tele.Context) *Message {
	var file string
	if p := c.Message().Photo; p != nil {
		file = p.FileID
	}
	return &Message{
		Group: getGroup(c),
		User:  getUser(c),
		Text:  c.Text(),
		File:  file,
		c:     c,
	}
}

func getGroup(c tele.Context) *Group {
	return &Group{ID: c.Chat().ID}
}

func getUser(c tele.Context) *User {
	return &User{ID: c.Sender().ID}
}
