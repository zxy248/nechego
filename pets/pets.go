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

func (g Gender) Icon() string {
	return genderData[g].Icon
}

func (g Gender) String() string {
	return genderData[g].Name
}

type Pet struct {
	Name    string
	Species Species
	Gender  Gender
	Birth   time.Time
}

func RandomPet(l float64) *Pet {
	return &Pet{
		Species: randomSpecies(l),
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
	var nameSuffix string
	if p.HasName() {
		nameSuffix = " "
	}
	return fmt.Sprintf("%s %s%s(%s)", p.Species, p.Name, nameSuffix, p.Gender.Icon())
}
