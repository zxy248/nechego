package command

import (
	"html"
	"strings"
)

var definitionPattern = "[^\\|]+"

func commandDefinition(s string) string {
	return html.EscapeString(strings.ToLower(strings.TrimSpace(s)))
}

func commandSubstitution(s string) string {
	return html.EscapeString(strings.TrimSpace(s))
}
