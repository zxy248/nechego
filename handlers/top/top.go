package top

import (
	"nechego/game"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

func trim(us []*game.User) []*game.User {
	const n = 5
	return us[:min(len(us), n)]
}

func whoFunc(c tele.Context) func(*game.User) string {
	return func(u *game.User) string {
		return tu.Link(c, u)
	}
}
