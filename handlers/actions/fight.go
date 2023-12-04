package actions

import (
	"math/rand"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	tu "nechego/teleutil"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Fight struct {
	Universe *game.Universe
}

var fightRe = handlers.Regexp("^!(драка|дуэль|поединок|атака|битва|схватка|сражение|бой|борьба)")

func (h *Fight) Match(c tele.Context) bool {
	return fightRe.MatchString(c.Text())
}

func (h *Fight) Handle(c tele.Context) error {
	reply, ok := tu.Reply(c)
	if !ok || c.Sender().ID == reply.ID {
		return c.Send(format.RepostMessage)
	}
	world, user1 := tu.Lock(c, h.Universe)
	user2 := world.User(reply.ID)
	defer world.Unlock()

	if time.Since(user2.LastMessage) > 10*time.Minute {
		return c.Send(format.NotOnline)
	}
	if !user1.Energy.Spend(0.33) {
		return c.Send(format.NoEnergy)
	}
	uw, ul, dr := game.Fight(user1, user2)

	s := []string{
		format.Fight(user1, user2),
		format.Win(tu.Link(c, uw), dr),
	}
	if rand.Float64() < 0.03 {
		if x, ok := moveRandomItem(uw.Inventory, ul.Inventory); ok {
			s = append(s, format.WinnerTook(tu.Link(c, uw), x))
		}
	}
	return c.Send(strings.Join(s, "\n"), tele.ModeHTML)
}

func moveRandomItem(dst, src *item.Set) (i *item.Item, ok bool) {
	i, ok = src.Random()
	if !ok {
		return nil, false
	}
	return i, src.Move(dst, i)
}
