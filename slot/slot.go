package slot

import "fmt"

type symbol int

const (
	bar symbol = iota
	grape
	lemon
	seven
)

func winning(v int) (s symbol, ok bool) {
	switch v {
	case 1:
		return bar, true
	case 22:
		return grape, true
	case 43:
		return lemon, true
	case 64:
		return seven, true
	}
	return 0, false
}

func multiplier(s symbol) int {
	switch s {
	case bar:
		return 4
	case grape:
		return 8
	case lemon:
		return 16
	case seven:
		return 32
	}
	panic(fmt.Sprintf("slot: unhandled symbol %v", s))
}

// Prize returns the amount of money you win with the given slot value
// and bet.
func Prize(value, bet int) (prize int, ok bool) {
	if s, ok := winning(value); ok {
		return bet * multiplier(s), true
	}
	return 0, false
}
