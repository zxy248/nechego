package pets

import (
	"fmt"
	"math/rand"
	"time"
)

type Gender int

const (
	Male Gender = iota
	Female

	numberOfGenders
)

func randomGender() Gender {
	return Gender(rand.Intn(int(numberOfGenders)))
}

func (g Gender) Emoji() string {
	return genderData[g].Emoji
}

func (g Gender) String() string {
	return genderData[g].Description
}

type Pet struct {
	Name    string
	Species Species
	Gender  Gender
	Birth   time.Time
}

func RandomPet(p float64) *Pet {
	return &Pet{
		Species: randomSpecies(p),
		Gender:  randomGender(),
		Birth:   time.Now(),
	}
}

func (p *Pet) HasName() bool {
	return p.Name != ""
}

func (p *Pet) Age() time.Duration {
	return time.Since(p.Birth)
}

func (p *Pet) String() string {
	space := ""
	if p.HasName() {
		space = " "
	}
	return fmt.Sprintf("%s %s%s(%s)", p.Species, p.Name, space, p.Gender.Emoji())
}
