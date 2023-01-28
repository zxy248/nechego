package pets

import (
	"fmt"
	"math/rand"
	"sort"
)

type Size int

const (
	Small Size = iota
	Medium
	Big
)

type Quality int

const (
	Common Quality = iota
	Rare
	Exotic
	Legendary
)

func (q Quality) String() string {
	switch q {
	case Common:
		return "обычный"
	case Rare:
		return "необычный"
	case Exotic:
		return "экзотический"
	case Legendary:
		return "легендарный"
	default:
		panic(fmt.Sprintf("unexpected quality %d", q))
	}
}

type Species int

const (
	Cat Species = iota
	Dog
	Hamster
	Rabbit
	Fox
	Bear
	Panda
	Koala
	Tiger
	Lion
	Cow
	Pig
	Frog
	Monkey
	Chicken
	Penguin
	Bird
	BabyChick
	Duck
	Eagle
	Owl
	Bat
	Wolf
	Boar
	Horse
	Unicorn
	Bee
	Worm
	Bug
	Butterfly
	Snail
	Beetle
	Ant
	Fly
	Mosquito
	Cricket
	Spider
	Scorpion
	Turtle
	Snake
	Lizard
	TRex
	Sauropod
	Octopus
	Squid
	Shrimp
	Crayfish
	Crab
	Dolphin
	Whale
	Shark
	Seal
	Crocodile
	Leopard
	Zebra
	Gorilla
	Mammon
	Elephant
	Camel
	Rhino
	Giraffe
	Kangaroo
	Ram
	Sheep
	Alpaca
	Goat
	Deer
	Rooster
	Turkey
	Peacock
	Swan
	Flamingo
	Hare
	Beaver
	Mouse
	Rat
	Chipmunk
	Hedgehog
	Parrot
	Dragon
	Caterpillar
)

func randomSpecies() Species {
	s := []Species{}
	r := rand.Float64()
	for k, v := range species {
		if r < v.Probability {
			s = append(s, k)
		}
	}
	if len(s) == 0 {
		panic("empty list")
	}

	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	sort.Slice(s, func(i, j int) bool {
		return species[s[i]].Probability < species[s[j]].Probability
	})
	few := 3
	if len(s) < few {
		few = len(s)
	}
	return s[rand.Intn(few)]
}

func (s Species) Emoji() string        { return species[s].Emoji }
func (s Species) String() string       { return species[s].Description }
func (s Species) Probability() float64 { return species[s].Probability }
func (s Species) Size() Size           { return species[s].Size }
func (s Species) Quality() Quality {
	p := s.Probability()
	switch {
	case p <= 0.01:
		return Legendary
	case p <= 0.05:
		return Exotic
	case p <= 0.20:
		return Rare
	default:
		return Common
	}
}
