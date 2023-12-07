package handlers

import (
	"fmt"
	"nechego/format"
	"nechego/game"
	"nechego/item"
	tu "nechego/teleutil"
	"nechego/valid"
	"strings"

	tele "gopkg.in/telebot.v3"
)

const InventoryCapacity = 20

func FullInventory(i *item.Set) bool {
	return i.Count() >= InventoryCapacity
}

func GetItems(s *item.Set, ks []int) []*item.Item {
	var items []*item.Item
	seen := map[*item.Item]bool{}
	for _, k := range ks {
		x, ok := s.ByKey(k)
		if !ok || seen[x] {
			break
		}
		seen[x] = true
		items = append(items, x)
	}
	return items
}

func MoveItems(dst, src *item.Set, items []*item.Item) (moved []*item.Item, bad *item.Item) {
	for _, x := range items {
		if !src.Move(dst, x) {
			return moved, x
		}
		moved = append(moved, x)
	}
	return
}

type Energy struct {
	Universe *game.Universe
}

var energyRe = Regexp("^!энергия")

func (h *Energy) Match(s string) bool {
	return energyRe.MatchString(s)
}

func (h *Energy) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	emoji := "🔋"
	if user.Energy < 0.5 {
		emoji = "🪫"
	}
	return c.Send(fmt.Sprintf("%s Запас энергии: %s",
		emoji, format.Energy(user.Energy)), tele.ModeHTML)
}

type NamePet struct {
	Universe *game.Universe
}

var namePetRe = Regexp("^!назвать (.+)")

func (h *NamePet) Match(s string) bool {
	return namePetRe.MatchString(s)
}

func (h *NamePet) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	pet, ok := user.Pet()
	if !ok {
		return c.Send("🐱 У вас нет питомца.")
	}

	e := pet.Species.Emoji()
	n := petName(c.Text())
	if n == "" {
		return c.Send(fmt.Sprintf("%s Такое имя не подходит для питомца.", e))
	}
	pet.Name = n
	s := fmt.Sprintf("%s Вы назвали питомца <code>%s</code>.", e, n)
	return c.Send(s, tele.ModeHTML)
}

func petName(s string) string {
	n := namePetRe.FindStringSubmatch(s)[1]
	if !valid.Name(n) {
		return ""
	}
	return strings.Title(n)
}
