package fishing

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	minimumWeight      = 0.05
	cheapThreshold     = 1000
	expensiveThreshold = 10000
)

var species = map[Species]struct {
	name          string
	normalWeight  float64
	maximumWeight float64
	constitution  Constitution
	pricePerKg    float64
	predator      bool
}{
	Pike:        {"Щука", 2.0, 35.0, Long, 400, true},
	Perch:       {"Окунь", 0.2, 5.0, Belly, 350, true},
	Zander:      {"Судак", 2.0, 18.0, Long, 400, true},
	Ruffe:       {"Ерш", 0.1, 0.2, Regular, 50, false},
	VolgaZander: {"Берш", 1.3, 2.0, Long, 430, true},
	Asp:         {"Жерех", 2.5, 4.5, Regular, 400, true},
	Chub:        {"Голавль", 0.75, 6.0, Regular, 300, false},
	Snakehead:   {"Змееголов", 3.0, 10.0, Long, 150, true},
	Burbot:      {"Налим", 4.5, 24.0, Long, 450, true},
	Eel:         {"Угорь", 2.0, 8.50, Long, 1500, true},
	Catfish:     {"Сом", 20.0, 150.0, Long, 500, true},
	Salmon:      {"Лосось", 4.0, 8.0, Regular, 1000, true},
	Grayling:    {"Хариус", 0.7, 1.4, Regular, 800, false},
	Trout:       {"Форель", 2.0, 10.0, Regular, 1000, true},
	Char:        {"Голец", 0.01, 0.025, Long, 50, false},
	Sturgeon:    {"Осетр", 18.0, 80.0, Long, 5000, true},
	Sterlet:     {"Стерлядь", 1.5, 8.0, Long, 1300, true},
	Carp:        {"Карп", 1.5, 24.0, Belly, 360, false},
	Goldfish:    {"Карась", 0.5, 5.0, Belly, 70, false},
	Tench:       {"Линь", 1.5, 7.5, Belly, 400, false},
	Bream:       {"Лещ", 1.0, 7.5, Belly, 100, false},
	Ide:         {"Язь", 1.0, 7.5, Regular, 300, false},
	Roach:       {"Плотва", 0.2, 2.0, Regular, 280, false},
	BigheadCarp: {"Толстолобик", 1.2, 16.0, Regular, 200, false},
	WhiteBream:  {"Белоглазка", 0.1, 0.8, Belly, 50, false},
	Rudd:        {"Красноперка", 0.3, 2.0, Belly, 100, false},
	Bleak:       {"Уклейка", 0.02, 0.06, Regular, 400, false},
	Nase:        {"Подуст", 0.4, 1.6, Regular, 180, false},
	Taimen:      {"Таймень", 4.0, 70.0, Long, 900, true},
}

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
		c *= 1.0
	case Belly:
		c *= math.Pi
	case Regular:
		c *= math.SqrtPi
	default:
		panic(fmt.Errorf("unexpected constitution %d", t))
	}
	return math.Pow(weight/c, 1.0/b)
}

func (t Constitution) randomLength(weight float64) float64 {
	mu := t.NormalLength(weight)
	sigma := mu / 10.0
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
)

func RandomSpecies() Species {
	return Species(rand.Intn(len(species)))
}

func (s Species) String() string             { return species[s].name }
func (s Species) NormalWeight() float64      { return species[s].normalWeight }
func (s Species) MaximumWeight() float64     { return species[s].maximumWeight }
func (s Species) Constitution() Constitution { return species[s].constitution }
func (s Species) PricePerKg() float64        { return species[s].pricePerKg }
func (s Species) Predator() bool             { return species[s].predator }

func (s Species) randomWeight() float64 {
	w := rand.NormFloat64()*s.weightStdDev() + s.NormalWeight()
	if w < minimumWeight {
		w = rand.Float64()*(s.NormalWeight()-minimumWeight) + minimumWeight
	}
	return w
}

func (s Species) weightStdDev() float64 {
	return (s.MaximumWeight() - s.NormalWeight()) / 3.0
}
