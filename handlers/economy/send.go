package economy

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	"nechego/money"
	tu "nechego/teleutil"
	"nechego/token"

	tele "gopkg.in/telebot.v3"
)

type Send struct {
	Universe *game.Universe
}

var sendRe = handlers.Regexp("^!(передать|отправить) ([0-9 ]+)")

func (h *Send) Match(c tele.Context) bool {
	return sendRe.MatchString(c.Text())
}

func (h *Send) Handle(c tele.Context) error {
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
	var sent []*item.Item
	keys := sendKeys(c.Text())
	for _, key := range keys {
		i, ok := user.Inventory.ByKey(key)
		if !ok {
			c.Send(format.BadKey(key), tele.ModeHTML)
			break
		}
		if !i.Transferable {
			c.Send(format.CannotTransfer(i), tele.ModeHTML)
			break
		}
		sendItem(user, target, i)
		sent = append(sent, i)
	}
	l1 := tu.Link(c, user)
	l2 := tu.Link(c, target)
	s := format.Transfered(l1, l2, sent)
	return c.Send(s, tele.ModeHTML)
}

func sendItem(from, to *game.User, i *item.Item) {
	from.Inventory.Remove(i)
	if c, ok := i.Value.(*money.Cash); ok {
		handlers.Pay(to, c.Money, from.Name)
		return
	}
	var x any
	if l, ok := i.Value.(*token.Letter); ok {
		x = l
	} else {
		x = &item.Box{From: from.Name, Content: i}
	}
	to.Mail.Add(item.New(x))
}

func sendKeys(s string) []int {
	m := sendRe.FindStringSubmatch(s)[2]
	return handlers.Numbers(m)
}
