package handlers

import (
	"fmt"
	"math/rand"
	"nechego/avatar"
	"nechego/fishing"
	"nechego/format"
	"nechego/game"
	"nechego/item"
	tu "nechego/teleutil"
	"nechego/valid"
	"strings"
	"time"
	"unicode/utf8"

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

type FishingRecords struct {
	Universe *game.Universe
}

var fishingRecordsRe = Regexp("^!рекорды")

func (h *FishingRecords) Match(s string) bool {
	return fishingRecordsRe.MatchString(s)
}

func (h *FishingRecords) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()
	byPrice := world.History.Top(fishing.Price, 10)
	byWeight := world.History.Top(fishing.Weight, 1)
	byLength := world.History.Top(fishing.Length, 1)
	for _, top := range [][]*fishing.Entry{byPrice, byWeight, byLength} {
		if len(top) == 0 {
			return c.Send(format.NoFishingRecords)
		}
	}
	return c.Send(format.FishingRecords(byPrice, byWeight[0], byLength[0]), tele.ModeHTML)
}

type Status struct {
	Universe  *game.Universe
	MaxLength int
}

var statusRe = Regexp("^!статус (.*)")

func (h *Status) Match(s string) bool {
	return statusRe.MatchString(s)
}

func (h *Status) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if reply, ok := tu.Reply(c); ok {
		// If the user has admin rights, they can set a status
		// for other users.
		if !user.Admin() {
			return c.Send("💬 Нельзя установить статус другому пользователю.")
		}
		user = world.User(reply.ID)
	}

	status := tu.Args(c, statusRe)[1]
	if utf8.RuneCountInString(status) > h.MaxLength {
		return c.Send(fmt.Sprintf("💬 Максимальная длина статуса %d символов.", h.MaxLength))
	}
	user.Status = status
	return c.Send("✅ Статус установлен.")
}

type Fight struct {
	Universe *game.Universe
}

var fightRe = Regexp("^!(драка|дуэль|поединок|атака|битва|схватка|сражение|бой|борьба)")

func (h *Fight) Match(s string) bool {
	return fightRe.MatchString(s)
}

func (h *Fight) Handle(c tele.Context) error {
	// Sanity check before locking the world.
	reply, ok := tu.Reply(c)
	if !ok {
		return c.Send(format.RepostMessage)
	}
	if c.Sender().ID == reply.ID {
		return c.Send(format.CannotAttackYourself)
	}

	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()
	opnt := world.User(reply.ID)

	// Can opponent fight back?
	if time.Since(opnt.LastMessage) > 10*time.Minute {
		return c.Send(format.NotOnline)
	}

	if !user.Energy.Spend(0.33) {
		return c.Send(format.NoEnergy)
	}

	// Fight begins.
	win, lose, elo := game.Fight(user, opnt)

	msg := format.NewConnector("\n\n")
	msg.Add(format.Fight(user, opnt))

	// The winner takes a random item.
	if rand.Float64() < 0.02 {
		if x, ok := moveRandomItem(win.Inventory, lose.Inventory); ok {
			msg.Add(format.WinnerTook(tu.Link(c, win), x))
		}
	}
	// The attacker drops a random item.
	if rand.Float64() < 0.04 {
		if x, ok := moveRandomItem(world.Floor, user.Inventory); ok {
			msg.Add(format.AttackerDrop(tu.Link(c, user), x))
		}
	}
	msg.Add(format.Win(tu.Link(c, win), elo))
	return c.Send(msg.String(), tele.ModeHTML)
}

func moveRandomItem(dst, src *item.Set) (i *item.Item, ok bool) {
	i, ok = src.Random()
	if !ok {
		return nil, false
	}
	return i, src.Move(dst, i)
}

type Profile struct {
	Universe *game.Universe
	Avatars  *avatar.Storage
}

var profileRe = Regexp("^!(профиль|стат)")

func (h *Profile) Match(s string) bool {
	return profileRe.MatchString(s)
}

func (h *Profile) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if u, ok := tu.Reply(c); ok {
		user = world.User(u.ID)
	}

	out := format.Profile(user)
	if a, ok := h.Avatars.Get(user.ID); ok {
		a.Caption = out
		return c.Send(a, tele.ModeHTML)
	}
	return c.Send(out, tele.ModeHTML)
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
