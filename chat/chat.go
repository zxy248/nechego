package chat

import tele "gopkg.in/telebot.v3"

type Photo struct {
	File string
}

type Group struct {
	ID int64
}

type User struct {
	ID int64
}

type Message struct {
	Group *Group
	User  *User
	Text  string
	File  string

	c tele.Context
}

func (m *Message) Reply(s string) error {
	return m.c.Send(s, tele.ModeHTML)
}

func (m *Message) ReplyPhoto(s string, p Photo) error {
	f := tele.File{FileID: p.File}
	return m.c.Send(&tele.Photo{File: f, Caption: s}, tele.ModeHTML)
}
