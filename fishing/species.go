package fishing

import (
	"math"
	"math/rand"
)

const (
	MinWeight          = 50e-3
	CheapThreshold     = 1e3
	ExpensiveThreshold = 1e4
)

type Constitution int

const (
	Long Constitution = iota
	Belly
	Regular
)

func (t Constitution) NormalLength(weight float64) float64 {
	c, b := 7.089, 3.096
	switch t {
	case Long:
		c *= 1.
	case Belly:
		c *= math.Pi
	case Regular:
		c *= math.SqrtPi
	default:
		panic("unknown constitution")
	}
	return math.Pow(weight/c, 1.0/b)
}

func (t Constitution) randomLength(weight float64) float64 {
	mu := t.NormalLength(weight)
	sigma := mu / math.Pow(math.Pi, 2)
	return rand.NormFloat64()*sigma + mu
}

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

	numberOfSpecies
)

func RandomSpecies() Species {
	return Species(rand.Intn(int(numberOfSpecies)))
}

func (s Species) String() string {
	return speciesData[s].name
}

func (s Species) NormalWeight() float64 {
	return speciesData[s].normalWeight
}

func (s Species) MaximumWeight() float64 {
	return speciesData[s].maximumWeight
}

func (s Species) Constitution() Constitution {
	return speciesData[s].constitution
}

func (s Species) PricePerKg() float64 {
	return speciesData[s].pricePerKg
}

func (s Species) Predator() bool {
	return speciesData[s].predator
}

func (s Species) randomWeight() float64 {
	weight := rand.NormFloat64()*s.weightSigma() + s.NormalWeight()
	if weight < MinWeight {
		weight = rand.Float64()*(s.NormalWeight()-MinWeight) + MinWeight
	}
	return weight
}

func (s Species) weightSigma() float64 {
	return (s.MaximumWeight() - s.NormalWeight()) / 3.0
}
