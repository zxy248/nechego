package format

import (
	"fmt"
	"nechego/game/reputation"
)

func ReputationTotal(mention string, score int) string {
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

func interpolatedReputationEmoji(score, lowest, highest int) string {
	diff := highest - lowest
	if diff == 0 {
		return "😐"
	}
	v := score - lowest
	x := float64(v) / float64(diff)

	emojis := [...]string{"👹", "👺", "👿", "😈", "🙂", "😌", "😊", "😇"}
	return emojis[int(x*float64(len(emojis)-1))]
}

func Reputation(r reputation.Reputation) string {
	return ReputationEmoji(r, "⭐️")
}

func ReputationEmoji(r reputation.Reputation, emoji string) string {
	return fmt.Sprintf("<code>%s %d</code>", emoji, r.Total())
}
