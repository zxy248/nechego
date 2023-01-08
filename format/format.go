package format

import (
	"fmt"
	"html"
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
