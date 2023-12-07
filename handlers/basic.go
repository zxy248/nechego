package handlers

import (
	"nechego/format"
	"nechego/game"
	tu "nechego/teleutil"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Help struct{}

var helpRe = Regexp("^!(помощь|команды|документ)")

func (h *Help) Match(c tele.Context) bool {
	return helpRe.MatchString(c.Text())
}

func (h *Help) Handle(c tele.Context) error {
	return c.Send("📖 <b>Документация:</b> nechego.pages.dev.", tele.ModeHTML)
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
