package command

import "strings"

var definitionPattern = "[^\\|]+"

func commandDefinition(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func commandSubstitution(s string) string {
	return strings.TrimSpace(s)
}
