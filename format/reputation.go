package format

import (
	"fmt"
	"nechego/game/reputation"
)

type Reputation struct {
	Score  int
	Factor float64
}

func (r Reputation) String(who string) string {
	return fmt.Sprintf("Репутация %s: %s", Name(who), r.rhsEmoji())
}

func (r Reputation) Updated(who string, d reputation.Direction) string {
	const format = "<b>⭐️ Репутация %s %s на <code>1</code></b>\n" +
		"Теперь репутация: %v"
	dd := "увеличена"
	if d == reputation.Down {
		dd = "понижена"
	}
	return fmt.Sprintf(format, Name(who), dd, r.rhsEmoji())
}

func (r Reputation) lhsEmoji() string {
	const format = "<code>%s %v</code>"
	return fmt.Sprintf(format, r.emoji(), r.Score)
}

func (r Reputation) rhsEmoji() string {
	const format = "<code>%v %s</code>"
	return fmt.Sprintf(format, r.Score, r.emoji())
}

func (r Reputation) emoji() string {
	e := [...]string{"👹", "👺", "👿", "😈", "😐", "🙂", "😌", "😊", "😇"}
	return e[int(r.Factor*float64(len(e)-1))]
}
