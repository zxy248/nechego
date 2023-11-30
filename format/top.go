package format

import (
	"fmt"
	"nechego/game"
)

func Top(head string, f func(*game.User) string, us []*game.User) string {
	c := NewConnector("\n")
	c.Add(head)
	for i, u := range us {
		c.Add(fmt.Sprintf("%s %s %s", Index(i), User(u), f(u)))
	}
	return c.String()
}

func TopRating(us []*game.User) string {
	f := func(u *game.User) string {
		return Rating(u.Rating)
	}
	return Top("<b>🏆 Боевой рейтинг</b>", f, us)
}

func TopRich(us []*game.User) string {
	f := func(u *game.User) string {
		return Money(u.Balance().Total())
	}
	return Top("💵 <b>Самые богатые пользователи</b>", f, us)
}

func TopStrength(us []*game.User) string {
	f := func(u *game.User) string {
		return Strength(u.Strength())
	}
	return Top("🏋️‍♀️ <b>Самые сильные пользователи</b>", f, us)
}

func TopPlain(title string, us []*game.User) string {
	head := fmt.Sprintf("<b>🏆 Топ %s</b>", title)
	f := func(u *game.User) string { return "" }
	return Top(head, f, us)
}
