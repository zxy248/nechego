package modifier

// Moder is implemented by any value that should have a corresponding modifier.
type Moder interface {
	Mod() (m *Mod, ok bool)
}

// Mod represents a modifier.
type Mod struct {
	Emoji       string
	Description string
	Multiplier  float64
}

// Mod trivially implements the Moder interface.
func (x *Mod) Mod() (m *Mod, ok bool) {
	return x, true
}

var (
	RatingFirst  = &Mod{"🥇", "Вы на 1-м месте в рейтинге.", +0.03}
	RatingSecond = &Mod{"🥈", "Вы на 2-м месте в рейтинге.", +0.02}
	RatingThird  = &Mod{"🥉", "Вы на 3-м месте в рейтинге.", +0.01}
)

// Set represents active modifiers.
type Set map[*Mod]bool

// Active is true if the given modifier is present in the set.
func (s Set) Active(m *Mod) bool {
	return s[m]
}

// Add adds the given modifier to the set.
func (s Set) Add(m *Mod) {
	s[m] = true
}

// List returns all modifiers from the set.
func (s Set) List() []*Mod {
	r := []*Mod{}
	for m := range s {
		r = append(r, m)
	}
	return r
}

// Sum returns the sum of all multipliers in the set.
func (s Set) Sum() float64 {
	r := 0.0
	for m := range s {
		r += m.Multiplier
	}
	return r
}
