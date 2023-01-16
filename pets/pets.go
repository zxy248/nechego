package pets

import (
	"fmt"
	"math/rand"
	"nechego/valid"
	"strings"
	"time"
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
	if !valid.Name(s) {
		return false
	}
	p.Name = strings.Title(s)
	return true
}
