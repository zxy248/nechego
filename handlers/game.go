package handlers

import (
	"errors"
	"fmt"
	"html"
	"nechego/game"
	"nechego/teleutil"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"
)

type Save struct {
	Universe *game.Universe
}

var saveRe = regexp.MustCompile("!сохран")

func (h *Save) Match(s string) bool {
	return saveRe.MatchString(s)
}

func (h *Save) Handle(c tele.Context) error {
	return h.Universe.SaveAll()
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
	u.GenerateHotkeys()
	items := u.ItemList()
	sort.Slice(items, func(i, j int) bool {
		return items[i].Hotkey < items[j].Hotkey
	})
	mention := teleutil.Mention(c, teleutil.Member(c, c.Sender()))
	head := fmt.Sprintf("<b>🗄 Инвентарь пользователя %s</b>", mention)
	lines := []string{head}
	for _, v := range items {
		i, ok := u.ItemByID(v.ItemID)
		if !ok {
			return errors.New("can't get item")
		}
		lines = append(lines, fmt.Sprintf("<code>%s:</code> %s", v.Hotkey, i.Value))
	}
	return c.Send(strings.Join(lines, "\n"), tele.ModeHTML)
}
