package fight

import (
	"errors"
	"math/rand"
	"nechego/model"
)

var ErrSameUser = errors.New("same user")

type StrengthFunc func(model.User) (float64, error)

type Settings struct {
	ChanceRatio  float64
	StrengthFunc StrengthFunc
}

type Fight struct {
	Attacker *Fighter
	Defender *Fighter
	settings Settings
}

func New(attacker, defender model.User, s Settings) (*Fight, error) {
	if attacker.ID == defender.ID {
		return nil, ErrSameUser
	}
	a, err := newFighter(attacker, s.StrengthFunc)
	if err != nil {
		return nil, err
	}
	d, err := newFighter(defender, s.StrengthFunc)
	if err != nil {
		return nil, err
	}
	f := &Fight{
		Attacker: a,
		Defender: d,
		settings: s,
	}
	return f, nil
}

func (f *Fight) Winner() *Fighter {
	a := f.Attacker.power(f.settings.ChanceRatio)
	d := f.Defender.power(f.settings.ChanceRatio)
	if a > d {
		return f.Attacker
	}
	return f.Defender
}

func (f *Fight) Loser() *Fighter {
	a := f.Attacker.power(f.settings.ChanceRatio)
	d := f.Defender.power(f.settings.ChanceRatio)
	if a > d {
		return f.Defender
	}
	return f.Attacker
}

type Fighter struct {
	model.User
	Strength float64
	R        float64
}

func newFighter(u model.User, f StrengthFunc) (*Fighter, error) {
	strength, err := f(u)
	if err != nil {
		return nil, err
	}
	return &Fighter{
		User:     u,
		R:        randomCoefficient(),
		Strength: strength,
	}, nil
}

func (f *Fighter) power(chanceRatio float64) float64 {
	return formula(f.Strength, f.R, chanceRatio)
}

type Outcome struct {
	Winner *Fighter
	Loser  *Fighter
}

func randomCoefficient() float64 {
	return rand.Float64()*2 - 1
}

func formula(strength, chance, chanceRatio float64) float64 {
	return (strength * (1 - chanceRatio)) + (strength * chance * chanceRatio)
}
