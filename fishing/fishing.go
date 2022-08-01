package fishing

import (
	"fmt"
	"math"
	"math/rand"
	"nechego/numbers"
)

type Constitution int

const (
	Long Constitution = iota
	Belly
	Regular
)

var MinFishWeight float64 = 0.05

type Species int

const (
	Pike Species = iota
	Perch
	Zander
	Ruffe
	VolgaZander
	Asp
	Chub
	Snakehead
	Burbot
	Eel
	Catfish
	Salmon
	Grayling
	Trout
	Char
	Sturgeon
	Sterlet
	Carp
	Goldfish
	Tench
	Bream
	Ide
	Roach
	BigheadCarp
	WhiteBream
	Rudd
	Bleak
	Nase
	Taimen
	NSpecies
)

func RandomSpecies() Species {
	return Species(rand.Intn(int(NSpecies)))
}

func (s Species) randomWeight() float64 {
	return numbers.UniNormal(
		MinFishWeight,
		speciesData[s].normalWeight,
		speciesData[s].maximumWeight)
}

func (s Species) String() string {
	return speciesData[s].name
}

func (s Species) Predator() bool {
	return speciesData[s].predator
}

func (s Species) NormalWeight() float64 {
	return speciesData[s].normalWeight
}

type Fish struct {
	Species
	Weight float64 // kilograms
	Length float64 // meters
}

func RandomFish() Fish {
	s := RandomSpecies()
	w := s.randomWeight()
	l := randomLength(w)
	return Fish{s, w, l}
}

func (f Fish) Price() int {
	return int(f.Weight * speciesData[f.Species].pricePerKg)
}

func (f Fish) Light() bool {
	return f.Weight < f.NormalWeight()
}

func (f Fish) Heavy() bool {
	return f.Weight >= f.NormalWeight()
}

func (f Fish) String() string {
	var length, weight string
	if f.Length < 1. {
		length = fmt.Sprintf("%.1f см", f.Length*100.)
	} else {
		length = fmt.Sprintf("%.2f м", f.Length)
	}
	if f.Weight < 1. {
		weight = fmt.Sprintf("%.1f г", f.Weight*1000.)
	} else {
		weight = fmt.Sprintf("%.2f кг", f.Weight)
	}
	return fmt.Sprintf("%s (%s, %s)", f.Species, weight, length)
}

func lengthFromWeight(w float64) float64 {
	c, b := 7.089, 3.096
	return math.Pow(w/c, 1.0/b)
}

func randomLength(w float64) float64 {
	l := lengthFromWeight(w)
	l += rand.NormFloat64() * l * (1.0 / 12)
	return l
}

type Fishes []Fish

func (f Fishes) Price() int {
	sum := 0
	for _, ff := range f {
		sum += ff.Price()
	}
	return sum
}

func (f Fishes) Weight() float64 {
	sum := 0.0
	for _, ff := range f {
		sum += ff.Weight
	}
	return sum
}

type Outcome int

const (
	Lost Outcome = iota
	Off
	Tear
	Seagrass
	Slip
	Release
	Collect
)

func (o Outcome) Success() bool {
	return o == goodOutcome()
}

func (o Outcome) String() string {
	return OutcomePrefix + outcomeDescriptions[o]
}

var OutcomePrefix = "🎣 "

func goodOutcome() Outcome {
	return Collect
}

var badOutcomes = []Outcome{
	Lost,
	Off,
	Tear,
	Seagrass,
	Slip,
	Release,
}

func badOutcome() Outcome {
	return badOutcomes[rand.Intn(len(badOutcomes))]
}

var SuccessChance = 0.5

type Session struct {
	Outcome
	Fish
}

func Cast() Session {
	return CastChance(SuccessChance)
}

func CastChance(win float64) Session {
	r := rand.Float64()
	outcome := badOutcome()
	if r <= win {
		outcome = goodOutcome()
	}
	fish := RandomFish()
	return Session{outcome, fish}
}
