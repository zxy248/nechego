package pets

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type Gender int

const (
	Male Gender = iota
	Female
)

func (g Gender) Emoji() string  { return genders[g].Emoji }
func (g Gender) String() string { return genders[g].Description }

type Pet struct {
	Name    string
	Species Species
	Gender  Gender
	Birth   time.Time
}

func Random() *Pet {
	return &Pet{
		Species: randomSpecies(),
		Gender:  Gender(rand.Intn(2)),
		Birth:   time.Now(),
	}
}

func (p *Pet) Age() time.Duration {
	return time.Since(p.Birth)
}

func (p *Pet) String() string {
	name := p.Name
	if name != "" {
		name = name + " "
	}
	return fmt.Sprintf("%s %s %s(%s)", p.Species.Emoji(),
		strings.Title(p.Species.String()), name, p.Gender.Emoji())
}

func (p *Pet) SetName(s string) bool {
	if !validName(s) {
		return false
	}
	p.Name = strings.Title(s)
	return true
}

func validName(s string) bool {
	if !utf8.ValidString(s) {
		return false
	}
	if utf8.RuneCountInString(s) > 40 {
		return false
	}
	for _, r := range s {
		if !unicode.Is(unicode.Cyrillic, r) && r != ' ' {
			return false
		}
	}
	return true
}
