package command

import (
	"fmt"
	"strings"
)

var (
	definitionPattern   = "[^\\|]+"
	separatorPattern    = "\\|"
	substitutionPattern = ".*"
	addPattern          = fmt.Sprintf(
		"^!добавить (%s)%s?(%s)",
		definitionPattern,
		separatorPattern,
		substitutionPattern,
	)
	removePattern = fmt.Sprintf(
		"^!(удалить|убрать) (%s)",
		definitionPattern,
	)
)

func sanitizeDefinition(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func sanitizeSubstitution(s string) string {
	return strings.TrimSpace(s)
}
