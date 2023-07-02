package format

import (
	"fmt"
	"nechego/game/reputation"
)

func Reputation(mention string, score int) string {
	return fmt.Sprintf("Репутация %s: %d", mention, score)
}

func ReputationUpdated(mention string, score int, d reputation.Dir) string {
	const format = "Репутация %s %s на 1\nТеперь репутация: %d"
	return fmt.Sprintf(format, mention, reputationDirectory(d), score)
}

func reputationDirectory(d reputation.Dir) string {
	switch d {
	case reputation.Up:
		return "увеличена"
	case reputation.Down:
		return "понижена"
	}
	panic(fmt.Sprintf("unknown reputation directory: %v", d))
}
