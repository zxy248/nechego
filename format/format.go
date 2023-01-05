package format

import (
	"fmt"
	"html"
)

func Mention(uid int64, name string) string {
	return fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, uid, html.EscapeString(name))
}
