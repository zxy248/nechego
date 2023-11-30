package format

import "fmt"

func List(title string, links []string) string {
	c := NewConnector("\n")
	c.Add(fmt.Sprintf("<b>📝 Список %s</b>", title))
	for _, l := range links {
		c.Add(fmt.Sprintf("<b>•</b> %s", l))
	}
	return c.String()
}
