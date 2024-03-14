package command

import "strings"

var definitionPattern = "[^\\|]+"

func sanitizeDefinition(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func sanitizeSubstitution(s string) string {
	return strings.TrimSpace(s)
}
