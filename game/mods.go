package game

import (
	"fmt"
	"nechego/fishing"
	"nechego/item"
	"nechego/pets"
	"nechego/token"
)

type Mod struct {
	Emoji       string
	Description string
	Multiplier  float64
}

func ModSum(ms []*Mod) float64 {
	x := 0.0
	for _, m := range ms {
		x += m.Multiplier
	}
	return x
}

var (
	RatingFirst  = &Mod{"🥇", "Вы на 1-м месте в рейтинге.", +0.03}
	RatingSecond = &Mod{"🥈", "Вы на 2-м месте в рейтинге.", +0.02}
	RatingThird  = &Mod{"🥉", "Вы на 3-м месте в рейтинге.", +0.01}
	TerribleLuck = &Mod{"☠️", "Вас преследуют неудачи.", -0.04}
	BadLuck      = &Mod{"🌧", "Вам не везёт.", -0.02}
	GreatLuck    = &Mod{"🍀", "Сегодня ваш день.", +0.02}
	GoodLuck     = &Mod{"🤞", "Вам везёт.", +0.04}
	LowEnergy    = &Mod{"😣", "Вы чувствуете себя уставшим.", -0.2}
	HighEnergy   = &Mod{"⚡️", "Вы полны сил.", 0.1}
	Rich         = &Mod{"🎩", "Вы богаты.", +0.05}
	Poor         = &Mod{"🗑️", "Вы бедны.", -0.05}
	Eblan        = &Mod{"😸", "Вы еблан.", -0.2}
	Admin        = &Mod{"👑", "Вы администратор.", 0.2}
	Pair         = &Mod{"💖", "У вас есть пара.", 0.1}
)

func luckMod(l float64) *Mod {
	switch {
	case l < 0.05:
		return TerribleLuck
	case l < 0.2:
		return BadLuck
	case l > 0.95:
		return GreatLuck
	case l > 0.8:
		return GoodLuck
	}
	return nil
}

func ratingMod(r int) *Mod {
	switch r {
	case 0:
		return RatingFirst
	case 1:
		return RatingSecond
	case 2:
		return RatingThird
	}
	return nil
}

func energyMod(e *Energy) *Mod {
	if e.Low() {
		return LowEnergy
	}
	if e.Full() {
		return HighEnergy
	}
	return nil
}

func moneyMod(b *Balance) *Mod {
	if b.Rich() {
		return Rich
	}
	if b.Poor() {
		return Poor
	}
	return nil
}

func petMod(p *pets.Pet) *Mod {
	var multiplier float64
	q := p.Species.Quality()
	switch q {
	case pets.Common:
		multiplier = 0.05
	case pets.Rare:
		multiplier = 0.10
	case pets.Exotic:
		multiplier = 0.15
	case pets.Legendary:
		multiplier = 0.20
	}
	pre := fmt.Sprintf("%s ", q)
	s := fmt.Sprintf("У вас есть %sпитомец: <code>%s</code>", pre, p)
	return &Mod{"🐱", s, multiplier}
}

func fishingRodMod(r *fishing.Rod) *Mod {
	return &Mod{"🎣", "Вы можете рыбачить.", 0.02 * float64(r.Level)}
}

func itemMod(i *item.Item) *Mod {
	switch x := i.Value.(type) {
	case *pets.Pet:
		return petMod(x)
	case *fishing.Rod:
		return fishingRodMod(x)
	case *token.Eblan:
		return Eblan
	case *token.Admin:
		return Admin
	case *token.Pair:
		return Pair
	}
	return nil
}

func (u *User) Mods() []*Mod {
	ms := []*Mod{
		luckMod(u.Luck()),
		ratingMod(u.RatingPosition),
		energyMod(&u.Energy),
		moneyMod(u.Balance()),
	}
	seen := map[item.Type]bool{}
	for _, x := range u.Inventory.List() {
		if seen[x.Type] {
			continue
		}
		seen[x.Type] = true
		if m := itemMod(x); m != nil {
			ms = append(ms, m)
		}
	}
	n := 0
	for _, m := range ms {
		if m != nil {
			ms[n] = m
			n++
		}
	}
	return ms[:n]
}
