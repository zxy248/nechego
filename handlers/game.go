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

type DailyEblan struct {
	Universe *game.Universe
}

var dailyEblanRe = regexp.MustCompile("^!еблан дня")

func (h *DailyEblan) Match(s string) bool {
	return dailyEblanRe.MatchString(s)
}

func (h *DailyEblan) Handle(c tele.Context) error {
	world := h.Universe.MustWorld(c.Chat().ID)
	world.Lock()
	defer world.Unlock()

	eblan, ok := world.DailyEblan()
	if !ok {
		return c.Send("😸")
	}
	mention := teleutil.Mention(c, teleutil.Member(c, tele.ChatID(eblan.TUID)))
	out := fmt.Sprintf("<b>Еблан дня</b> — %s 😸", mention)
	return c.Send(out, tele.ModeHTML)
}

type DailyAdmin struct {
	Universe *game.Universe
}

var dailyAdminRe = regexp.MustCompile("^!админ дня")

func (h *DailyAdmin) Match(s string) bool {
	return dailyAdminRe.MatchString(s)
}

func (h *DailyAdmin) Handle(c tele.Context) error {
	world := h.Universe.MustWorld(c.Chat().ID)
	world.Lock()
	defer world.Unlock()

	admin, ok := world.DailyAdmin()
	if !ok {
		return c.Send("👑")
	}
	m := teleutil.Mention(c, teleutil.Member(c, tele.ChatID(admin.TUID)))
	out := fmt.Sprintf("<b>Админ дня</b> — %s 👑", m)
	return c.Send(out, tele.ModeHTML)
}

type DailyPair struct {
	Universe *game.Universe
}

var dailyPairRe = regexp.MustCompile("^!пара дня")

func (h *DailyPair) Match(s string) bool {
	return dailyPairRe.MatchString(s)
}

func (h *DailyPair) Handle(c tele.Context) error {
	world := h.Universe.MustWorld(c.Chat().ID)
	world.Lock()
	defer world.Unlock()

	pair, ok := world.DailyPair()
	if !ok {
		return c.Send("💔")
	}
	m0 := teleutil.Mention(c, teleutil.Member(c, tele.ChatID(pair[0].TUID)))
	m1 := teleutil.Mention(c, teleutil.Member(c, tele.ChatID(pair[1].TUID)))
	out := fmt.Sprintf("<b>✨ Пара дня</b> — %s 💘 %s", m0, m1)
	return c.Send(out, tele.ModeHTML)
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
