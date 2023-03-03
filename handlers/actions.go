package handlers

import (
	"nechego/format"
	"nechego/game"
	"nechego/item"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Friends struct {
	Universe *game.Universe
}

var friendsRe = re("^!(друзья|друж)")

func (h *Friends) Match(s string) bool {
	return friendsRe.MatchString(s)
}

func (h *Friends) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if reply, ok := tu.Reply(c); ok {
		if user.TUID == reply.ID {
			return c.Send(format.CannotFriend)
		}
		target := world.UserByID(reply.ID)
		if user.Friends.With(target) {
			user.Friends.Remove(target)
			return c.Send(format.FriendRemoved(
				tu.Mention(c, user), tu.Mention(c, target)),
				tele.ModeHTML)
		} else {
			user.Friends.Add(target)
			if game.MutualFriends(user, target) {
				return c.Send(format.MutualFriends(
					tu.Mention(c, user), tu.Mention(c, target)),
					tele.ModeHTML)
			} else {
				return c.Send(format.FriendAdded(
					tu.Mention(c, user), tu.Mention(c, target)),
					tele.ModeHTML)
			}
		}
	}
	list := user.Friends.List()
	friends := make([]format.Friend, 0, len(list))
	for _, id := range list {
		target := world.UserByID(id)
		friends = append(friends, format.Friend{
			Mention: tu.Mention(c, target),
			Mutual:  game.MutualFriends(user, target),
		})
	}
	return c.Send(format.FriendList(tu.Mention(c, user), friends), tele.ModeHTML)
}

type Transfer struct {
	Universe *game.Universe
}

var transferRe = re("^!(передать|отправить) (.*)")

func (h *Transfer) Match(s string) bool {
	return transferRe.MatchString(s)
}

func (h *Transfer) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	reply, ok := tu.Reply(c)
	if !ok {
		return c.Send(format.RepostMessage)
	}

	target := world.UserByID(reply.ID)
	if !target.Friends.With(user) {
		return c.Send(format.NonFriendTransfer)
	}

	transfered := []*item.Item{}
	for _, key := range tu.NumArg(c, transferRe, 2) {
		item, ok := user.Inventory.ByKey(key)
		if !ok {
			c.Send(format.BadKey(key), tele.ModeHTML)
			break
		}
		if !user.Transfer(target, item) {
			c.Send(format.CannotTransfer(item), tele.ModeHTML)
			break
		}
		transfered = append(transfered, item)
	}
	return c.Send(format.Transfered(
		tu.Mention(c, user), tu.Mention(c, target), transfered...),
		tele.ModeHTML)
}
