package command

import (
	"fmt"
	"strings"

	"github.com/zxy248/nechego/handlers"
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
	cfgPattern = handlers.NewRegexp("^!конфиг(.*)")
)

func sanitizeDefinition(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func sanitizeSubstitution(s string) string {
	return strings.TrimSpace(s)
}
