package format

import (
	"fmt"
	"html"
	"nechego/fishing"
	"nechego/game"
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
