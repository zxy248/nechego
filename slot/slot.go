package slot

type symbol int

const (
	zero symbol = iota
	bar
	grape
	lemon
	seven
)

var lines = map[int]symbol{
	1:  bar,
	22: grape,
	43: lemon,
	64: seven,
}

var multipliers = map[symbol]int{
	bar:   4,
	grape: 8,
	lemon: 16,
	seven: 32,
}

func line(v int) symbol {
	if s, ok := lines[v]; ok {
		return s
	}
	return zero
}

func multiplier(s symbol) int {
	return multipliers[s]
}

// Prize returns the amount of money you win with the given slot value
// and bet.
func Prize(value, bet int) int {
	return bet * multiplier(line(value))
}
