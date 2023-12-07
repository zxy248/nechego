package handlers

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers/parse"
	"nechego/item"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

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

	target := world.User(reply.ID)
	if !target.Friends.With(user.ID) {
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
