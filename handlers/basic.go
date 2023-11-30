package handlers

import (
	"errors"
	"fmt"
	"math/rand"
	"nechego/avatar"
	"nechego/format"
	"nechego/game"
	"nechego/handlers/parse"
	tu "nechego/teleutil"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Help struct{}

var helpRe = Regexp("^!(помощь|команды|документ)")

func (h *Help) Match(s string) bool {
	return helpRe.MatchString(s)
}

func (h *Help) Handle(c tele.Context) error {
	return c.Send("📖 <b>Документация:</b> nechego.pages.dev.", tele.ModeHTML)
}

type Top struct {
	Universe *game.Universe
}

func (h *Top) Match(s string) bool {
	_, ok := topCommand(s)
	return ok
}

func (h *Top) Handle(c tele.Context) error {
	text, _ := topCommand(c.Text())
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	us := world.RandomUserIDs(3 + rand.Intn(3))
	s := []string{fmt.Sprintf("<b>🏆 Топ %s</b>", text)}
	for i, u := range us {
		s = append(s, fmt.Sprintf("<i>%d.</i> %s", 1+i, tu.Link(c, u)))
	}
	return c.Send(strings.Join(s, "\n"), tele.ModeHTML)
}

func topCommand(s string) (text string, ok bool) {
	return textCommand(parse.Match("!топ"), s)
}

type Avatar struct {
	Universe *game.Universe
	Avatars  *avatar.Storage
}

var avatarRe = Regexp("^!ава")

func (h *Avatar) Match(s string) bool {
	return avatarRe.MatchString(s)
}

func (h *Avatar) Handle(c tele.Context) error {
	target := c.Sender().ID
	photo := c.Message().Photo
	if reply, ok := tu.Reply(c); ok {
		// If the user has admin rights, they can set an
		// avatar for other users.
		world, user := tu.Lock(c, h.Universe)
		admin := user.Admin()
		world.Unlock()
		if !admin {
			return c.Send("📷 Нельзя установить аватар другому пользователю.")
		}
		target = reply.ID
	}

	if photo == nil {
		if avatar, ok := h.Avatars.Get(target); ok {
			return c.Send(avatar)
		}
		return c.Send("📷 Прикрепите изображение.")
	}
	if err := h.Avatars.Set(target, photo); errors.Is(err, avatar.ErrSize) {
		return c.Send("📷 Максимальный размер аватара %dx%d пикселей.",
			h.Avatars.MaxWidth, h.Avatars.MaxHeight)
	} else if err != nil {
		return err
	}
	return c.Send("📸 Аватар установлен.")
}

type Ban struct {
	Universe   *game.Universe
	DurationHr int // Ban duration in hours.
}

var banRe = Regexp("^!бан")

func (h *Ban) Match(s string) bool {
	return banRe.MatchString(s)
}

func (h *Ban) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if !user.Admin() {
		return c.Send(format.AdminsOnly)
	}
	reply, ok := tu.Reply(c)
	if !ok {
		return c.Send(format.RepostMessage)
	}
	target := world.User(reply.ID)
	duration := time.Hour * time.Duration(h.DurationHr)
	target.BannedUntil = time.Now().Add(duration)
	return c.Send(format.UserBanned(h.DurationHr))
}

type Unban struct {
	Universe *game.Universe
}

var unbanRe = Regexp("^!разбан")

func (h *Unban) Match(s string) bool {
	return unbanRe.MatchString(s)
}

func (h *Unban) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if !user.Admin() {
		return c.Send(format.AdminsOnly)
	}
	reply, ok := tu.Reply(c)
	if !ok {
		return c.Send(format.RepostMessage)
	}
	world.User(reply.ID).BannedUntil = time.Time{}
	return c.Send(format.UserUnbanned)
}
