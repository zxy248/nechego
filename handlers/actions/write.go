package actions

import (
	"fmt"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	tu "nechego/teleutil"
	"nechego/token"

	tele "gopkg.in/telebot.v3"
)

type Write struct {
	Universe *game.Universe
}

var writeRe = handlers.Regexp("(?s)^!написать (письмо )?(.+)")

func (h *Write) Match(c tele.Context) bool {
	return writeRe.MatchString(c.Text())
}

func (h *Write) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if handlers.FullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}
	l := &token.Letter{
		Author: user.Name,
		Text:   letterText(c.Text()),
	}
	i := item.New(l)
	user.Inventory.Add(i)
	m := format.User(user)
	s := fmt.Sprintf("✍️ %s пишет письмо.", m)
	return c.Send(s, tele.ModeHTML)
}

func letterText(s string) string {
	return writeRe.FindStringSubmatch(s)[2]
}
