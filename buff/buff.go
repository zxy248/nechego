package buff

import (
	"nechego/modifier"
	"time"
)

// Buff represents an applied effect.
type Buff int

const (
	Beer Buff = iota
)

var data = map[Buff]struct {
	emoji      string
	string     string
	multiplier float64
}{
	Beer: {"🍻", "Вы выпили пиво.", 0.25},
}

func (b Buff) Emoji() string  { return data[b].emoji }
func (b Buff) String() string { return data[b].string }

// Mod implements the modifier.Moder interface.
func (b Buff) Mod() (*modifier.Mod, bool) {
	return &modifier.Mod{
		Emoji:       b.Emoji(),
		Description: b.String(),
		Multiplier:  data[b].multiplier,
	}, true
}

// Set is a map from a Buff to the end of its effect's duration.
type Set map[Buff]time.Time

// Active returns true if the Buff b is currently active.
func (s Set) Active(b Buff) bool {
	return s[b].After(time.Now())
}

// Apply activates the Buff b for the given duration.
func (s Set) Apply(b Buff, d time.Duration) {
	s[b] = time.Now().Add(d)
}

// List returns all active Buffs.
func (s Set) List() []Buff {
	r := []Buff{}
	for b := range s {
		if s.Active(b) {
			r = append(r, b)
		}
	}
	return r
}
