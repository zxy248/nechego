package actions

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Transfer struct {
	Universe *game.Universe
}

var transferRe = handlers.Regexp("^!(передать|отправить) ([0-9 ]+)")

func (h *Transfer) Match(c tele.Context) bool {
	return transferRe.MatchString(c.Text())
}

func (h *Transfer) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	reply, ok := tu.Reply(c)
	if !ok {
		return c.Send(format.RepostMessage)
	}
	target := world.User(reply.ID)
	if !target.Friends.With(user.ID) {
		return c.Send(format.NonFriendTransfer)
	}
	var transfered []*item.Item
	keys := transferKeys(c.Text())
	for _, key := range keys {
		i, ok := user.Inventory.ByKey(key)
		if !ok {
			c.Send(format.BadKey(key), tele.ModeHTML)
			break
		}
		if !user.Transfer(target, i) {
			c.Send(format.CannotTransfer(i), tele.ModeHTML)
			break
		}
		transfered = append(transfered, i)
	}
	l1 := tu.Link(c, user)
	l2 := tu.Link(c, target)
	s := format.Transfered(l1, l2, transfered)
	return c.Send(s, tele.ModeHTML)
}

func transferKeys(s string) []int {
	m := transferRe.FindStringSubmatch(s)[2]
	return handlers.Numbers(m)
}
