package actions

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Friends struct {
	Universe *game.Universe
}

var friendsRe = handlers.Regexp("^!(друзья|друж)")

func (h *Friends) Match(c tele.Context) bool {
	return friendsRe.MatchString(c.Text())
}

func (h *Friends) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	var s string
	if reply, ok := tu.Reply(c); ok {
		target := world.User(reply.ID)
		s = h.handleReply(c, world, user, target)
	} else {
		s = h.handleList(c, world, user)
	}
	return c.Send(s, tele.ModeHTML)
}

func (h *Friends) handleReply(c tele.Context, w *game.World, u1, u2 *game.User) string {
	l1 := tu.Link(c, u1)
	l2 := tu.Link(c, u2)
	if u1.Friends.With(u2.ID) {
		u1.Friends.Remove(u2.ID)
		return format.FriendRemoved(l1, l2)
	}
	u1.Friends.Add(u2.ID)
	if u1.MutualFriends(u2) {
		return format.MutualFriends(l1, l2)
	}
	return format.FriendAdded(l1, l2)
}

func (*Friends) handleList(c tele.Context, w *game.World, u *game.User) string {
	list := u.Friends.List()
	friends := map[string]bool{}
	for _, id := range list {
		v := w.User(id)
		l := tu.Link(c, v)
		m := u.MutualFriends(v)
		friends[l] = m
	}
	l := tu.Link(c, u)
	return format.FriendList(l, friends)
}
