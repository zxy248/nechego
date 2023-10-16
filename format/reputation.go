package format

import (
	"fmt"
	"nechego/game/reputation"
)

func ReputationScore(mention string, x string) string {
	return fmt.Sprintf("Репутация %s: %s", Name(mention), x)
}

func ReputationUpdated(mention string, score string, d reputation.Direction) string {
	return Lines(
		Bold(Words("⭐️", "Репутация", Name(mention),
			reputationDirection(d), "на", Code("1"))),
		Words("Теперь репутация:", score),
	)
}

func ReputationPrefix(score, low, high int) string {
	emoji := reputationEmoji(score, low, high)
	return Code(Words(emoji, Value(score)))
}

func ReputationSuffix(score, low, high int) string {
	emoji := reputationEmoji(score, low, high)
	return Code(Words(Value(score), emoji))
}

func reputationDirection(d reputation.Direction) string {
	switch d {
	case reputation.Up:
		return "увеличена"
	case reputation.Down:
		return "понижена"
	}
	panic(fmt.Sprintf("unknown reputation directory: %v", d))
}

func reputationEmoji(score, low, high int) string {
	diff := high - low
	if diff == 0 {
		return "😐"
	}
	v := score - low
	x := float64(v) / float64(diff)

	emojis := [...]string{"👹", "👺", "👿", "😈", "🙂", "😌", "😊", "😇"}
	return emojis[int(x*float64(len(emojis)-1))]
}
