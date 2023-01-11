package format

import (
	"fmt"
	"html"
	"nechego/fishing"
	"nechego/game"
	"nechego/modifier"
	"strings"
)

const empty = "<code>. . .</code>"

func Mention(uid int64, name string) string {
	return fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, uid, html.EscapeString(name))
}

func Items(items []*game.Item) []string {
	lines := []string{}
	for i, v := range items {
		lines = append(lines, fmt.Sprintf("<code>%v ≡ </code> %s", i, Item(v)))
	}
	if len(lines) == 0 {
		return []string{empty}
	}
	return lines
}

func Item(i *game.Item) string {
	return fmt.Sprintf("<code>%s</code>", i.Value)
}

func Products(products []*game.Product) []string {
	lines := []string{}
	for i, p := range products {
		line := fmt.Sprintf("<code>%v ≡ </code> %s, %s", i, Item(p.Item), Money(p.Price))
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return []string{empty}
	}
	return lines
}

func Money(q int) string {
	return fmt.Sprintf("<code>%d ₽</code>", q)
}

func Energy(e int) string {
	return fmt.Sprintf("<code>%d ⚡</code>", e)
}

func Fish(f *fishing.Fish) string {
	return fmt.Sprintf("<code>%s</code>", f)
}

func Rating(r float64) string {
	return fmt.Sprintf("<code>%.1f ⚜️</code>", r)
}

func Strength(s float64) string {
	return fmt.Sprintf("<code>%.2f 💪</code>", s)
}

func Messages(n int) string {
	return fmt.Sprintf("<code>%d ✉️</code>", n)
}

func Status(s string) string {
	return fmt.Sprintf("<i>%s</i>", s)
}

func Key(k int) string {
	return fmt.Sprintf("<code>#%d</code>", k)
}

func ModifierEmojis(m []*modifier.Mod) string {
	emojis := []string{}
	for _, x := range m {
		if x.Emoji != "" {
			emojis = append(emojis, x.Emoji)
		}
	}
	return fmt.Sprintf("<code>%s</code>", strings.Join(emojis, "·"))
}

func ModifierDescriptions(m []*modifier.Mod) string {
	descs := []string{}
	for _, x := range m {
		descs = append(descs, x.Description)
	}
	return fmt.Sprintf("<i>%s</i>", strings.Join(descs, "\n"))
}

func ModifierTitles(m []*modifier.Mod) string {
	titles := []string{}
	for _, x := range m {
		if x.Title != "" {
			titles = append(titles, x.Title)
		}
	}
	if len(titles) == 0 {
		titles = append(titles, "пользователь")
	}
	titles[0] = strings.Title(titles[0])
	return strings.Join(titles, " ")
}
