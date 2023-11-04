package handlers

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers/parse"
	"nechego/item"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Friends struct {
	Universe *game.Universe
}

var friendsRe = Regexp("^!(друзья|друж)")

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
				tu.Link(c, user), tu.Link(c, target)),
				tele.ModeHTML)
		} else {
			user.Friends.Add(target)
			if game.MutualFriends(user, target) {
				return c.Send(format.MutualFriends(
					tu.Link(c, user), tu.Link(c, target)),
					tele.ModeHTML)
			} else {
				return c.Send(format.FriendAdded(
					tu.Link(c, user), tu.Link(c, target)),
					tele.ModeHTML)
			}
		}
	}
	list := user.Friends.List()
	friends := make([]format.Friend, 0, len(list))
	for _, id := range list {
		target := world.UserByID(id)
		friends = append(friends, format.Friend{
			Who:    tu.Link(c, target),
			Mutual: game.MutualFriends(user, target),
		})
	}
	return c.Send(format.FriendList(tu.Link(c, user), friends), tele.ModeHTML)
}

type Transfer struct {
	Universe *game.Universe
}

func (h *Transfer) Match(s string) bool {
	_, ok := transferCommand(s)
	return ok
}

func (h *Transfer) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()
	keys, _ := transferCommand(c.Text())

	reply, ok := tu.Reply(c)
	if !ok {
		return c.Send(format.RepostMessage)
	}

	target := world.UserByID(reply.ID)
	if !target.Friends.With(user) {
		return c.Send(format.NonFriendTransfer)
	}

	transfered := []*item.Item{}
	for _, key := range keys {
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
		tu.Link(c, user), tu.Link(c, target), transfered...),
		tele.ModeHTML)
}

func transferCommand(s string) (keys []int, ok bool) {
	return numCommand(parse.Match("!передать", "!отправить"), s)
}
