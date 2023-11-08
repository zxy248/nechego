package top

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"
	"regexp"

	tele "gopkg.in/telebot.v3"
)

type Top struct {
	Universe *game.Universe
	Regexp   *regexp.Regexp
	Sort     game.UserSortFunc
	Format   func([]*game.User) string
}

func (h *Top) Match(c tele.Context) bool {
	return h.Regexp.MatchString(c.Text())
}

func (h *Top) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	us := world.SortedUsers(h.Sort)
	us = us[:min(len(us), 5)]
	return c.Send(h.Format(us), tele.ModeHTML)
}

func Rating(u *game.Universe) *Top {
	return &Top{
		Universe: u,
		Regexp:   handlers.Regexp("^!(рейтинг|ммр|эло)"),
		Sort:     game.ByElo,
		Format:   format.TopRating,
	}
}

func Rich(u *game.Universe) *Top {
	return &Top{
		Universe: u,
		Regexp:   handlers.Regexp("^!топ бога[тч]"),
		Sort:     game.ByWealth,
		Format:   format.TopRich,
	}
}

func Strength(u *game.Universe) *Top {
	return &Top{
		Universe: u,
		Regexp:   handlers.Regexp("^!топ сил"),
		Sort:     game.ByStrength,
		Format:   format.TopStrength,
	}
}
