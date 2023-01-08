package handlers

import (
	"errors"
	"fmt"
	"html"
	"nechego/format"
	"nechego/game"
	"nechego/teleutil"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"
)

type Save struct {
	Universe *game.Universe
}

var saveRe = regexp.MustCompile("^!сохран")

func (h *Save) Match(s string) bool {
	return saveRe.MatchString(s)
}

func (h *Save) Handle(c tele.Context) error {
	if err := h.Universe.SaveAll(); err != nil {
		return err
	}
	return c.Send("💾 Игра сохранена.")
}

type Name struct{}

var nameRe = regexp.MustCompile("!имя (.*)")

func (h *Name) Match(s string) bool {
	return nameRe.MatchString(s)
}

func (h *Name) Handle(c tele.Context) error {
	name := html.EscapeString(teleutil.Args(c, nameRe)[1])
	const max = 16
	if utf8.RuneCountInString(name) > max {
		return c.Send(fmt.Sprintf("⚠️ Максимальная длина имени %d символов.", max))
	}

	user := c.Sender()
	if err := teleutil.Promote(c, teleutil.Member(c, user)); err != nil {
		return err
	}
	if err := c.Bot().SetAdminTitle(c.Chat(), user, name); err != nil {
		return err
	}
	return c.Send(fmt.Sprintf("Имя <b>%s</b> установлено ✅", name), tele.ModeHTML)
}

type Inventory struct {
	Universe *game.Universe
}

var inventoryRe = regexp.MustCompile("^!инвентарь")

func (h *Inventory) Match(s string) bool {
	return inventoryRe.MatchString(s)
}

func (h *Inventory) Handle(c tele.Context) error {
	world := h.Universe.MustWorld(c.Chat().ID)
	world.Lock()
	defer world.Unlock()

	u, ok := world.UserByID(c.Sender().ID)
	if !ok {
		return errors.New("user not found")
	}
	items := u.ListInventory()
	mention := teleutil.Mention(c, teleutil.Member(c, c.Sender()))
	head := fmt.Sprintf("<b>🗄 Инвентарь пользователя %s</b>", mention)
	lines := append([]string{head}, format.Items(items)...)
	return c.Send(strings.Join(lines, "\n"), tele.ModeHTML)
}

type Drop struct {
	Universe *game.Universe
}

var dropRe = regexp.MustCompile("^!выкинуть (.*)")

func (h *Drop) Match(s string) bool {
	return dropRe.MatchString(s)
}

func (h *Drop) Handle(c tele.Context) error {
	world := h.Universe.MustWorld(c.Chat().ID)
	world.Lock()
	defer world.Unlock()

	user, ok := world.UserByID(c.Sender().ID)
	if !ok {
		return errors.New("user not found")
	}
	k, err := strconv.Atoi(teleutil.Args(c, dropRe)[1])
	if err != nil {
		return c.Send("#⃣ Укажите номер предмета.")
	}
	item, ok := user.ItemByHotkey(k)
	if !ok {
		return c.Send("🗄 Такого предмета нет в инвентаре.")
	}
	if ok := world.Drop(user, item); !ok {
		return c.Send("♻ Вы не можете выбросить этот предмет.")
	}
	out := fmt.Sprintf("🚮 Вы выбросили %s.", format.Item(item))
	return c.Send(out, tele.ModeHTML)
}

type Floor struct {
	Universe *game.Universe
}

var floorRe = regexp.MustCompile("^!пол")

func (h *Floor) Match(s string) bool {
	return floorRe.MatchString(s)
}

func (h *Floor) Handle(c tele.Context) error {
	world := h.Universe.MustWorld(c.Chat().ID)
	world.Lock()
	defer world.Unlock()

	items := world.ListFloor()
	head := fmt.Sprintf("<b>🗑 Пол</b>")
	lines := append([]string{head}, format.Items(items)...)
	return c.Send(strings.Join(lines, "\n"), tele.ModeHTML)
}
