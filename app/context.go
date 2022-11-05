package app

import (
	"nechego/input"
	"nechego/model"

	tele "gopkg.in/telebot.v3"
)

const (
	messageKey   = "message"
	groupKey     = "group"
	userKey      = "user"
	replyUserKey = "replyUser"
)

func addMessage(c tele.Context, m *input.Message) tele.Context {
	c.Set(messageKey, m)
	return c
}

func getMessage(c tele.Context) *input.Message {
	return c.Get(messageKey).(*input.Message)
}

func addGroup(c tele.Context, g model.Group) tele.Context {
	c.Set(groupKey, g)
	return c
}

func getGroup(c tele.Context) model.Group {
	return c.Get(groupKey).(model.Group)
}

func addUser(c tele.Context, u model.User) tele.Context {
	c.Set(userKey, u)
	return c
}

func getUser(c tele.Context) model.User {
	return c.Get(userKey).(model.User)
}

func addReplyUser(c tele.Context, u model.User) tele.Context {
	c.Set(replyUserKey, u)
	return c
}

func getReplyUser(c tele.Context) model.User {
	u, _ := maybeGetReplyUser(c)
	return u
}

func maybeGetReplyUser(c tele.Context) (model.User, bool) {
	u, ok := c.Get(replyUserKey).(model.User)
	return u, ok
}
