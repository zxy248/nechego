package pets

import (
	"math/rand"
	"sort"
)

type Quality int

const (
	Common Quality = iota
	Rare
	Exotic
	Legendary

	numberOfQualities
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
		panic("unknown quality")
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

	numberOfSpecies
)

func randomSpecies(l float64) Species {
	s := []Species{}
	for k, v := range speciesData {
		if l < v.Rarity {
			s = append(s, k)
		}
	}
	if len(s) == 0 {
		panic("no species")
	}

	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	sort.Slice(s, func(i, j int) bool {
		return speciesData[s[i]].Rarity < speciesData[s[j]].Rarity
	})

	few := 4
	if len(s) < few {
		few = len(s)
	}
	return s[rand.Intn(few)]
}

func (s Species) Rarity() float64 {
	return speciesData[s].Rarity
}

func (s Species) RarityMultiplier() float64 {
	return 1.0 / speciesData[s].Rarity
}

func (s Species) Icon() string {
	return speciesData[s].Icon
}

func (s Species) String() string {
	return speciesData[s].Name
}

func (s Species) Quality() Quality {
	switch r := s.Rarity(); {
	case r < 0.01:
		return Legendary
	case r < 0.05:
		return Exotic
	case r < 0.20:
		return Rare
	default:
		return Common
	}
}
