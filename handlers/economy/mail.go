package economy

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	"nechego/money"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Mail struct {
	Universe *game.Universe
}

var mailRe = handlers.Regexp("^!(почта|зарплата|средства|получить)")

func (h *Mail) Match(c tele.Context) bool {
	return mailRe.MatchString(c.Text())
}

func (h *Mail) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if handlers.FullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}
	var got []*item.Item
	for _, i := range user.Mail.List() {
		receiveItem(user, i)
		got = append(got, i)
	}
	l := tu.Link(c, user)
	s := format.Mail(l, got)
	return c.Send(s, tele.ModeHTML)
}

func receiveItem(u *game.User, i *item.Item) {
	u.Mail.Remove(i)
	if t, ok := i.Value.(*money.Transfer); ok {
		i = item.New(&money.Cash{Money: t.Money})
	}
	u.Inventory.Add(i)
}
