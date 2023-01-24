package game

import (
	"fmt"
	"nechego/modifier"
	"nechego/pets"
)

func (u *User) Modset(w *World) modifier.Set {
	set := modifier.Set{}
	if u.Admin() {
		set.Add(modifier.Admin)
	}
	if u.Eblan() {
		set.Add(modifier.Eblan)
	}
	if u.Energy == 0 {
		set.Add(modifier.NoEnergy)
	}
	if u.Energy == EnergyCap {
		set.Add(modifier.FullEnergy)
	}
	if u.Energy > EnergyCap {
		set.Add(modifier.MuchEnergy)
	}
	if u.Rich() {
		set.Add(modifier.Rich)
	}
	if u.Poor() {
		set.Add(modifier.Poor)
	}
	if u.InDebt() {
		set.Add(modifier.Debtor)
	}
	if u.Inventory.Count() > InventorySize {
		set.Add(modifier.Heavy)
	}
	if l, ok := luckModifier(u.Luck()); ok {
		set.Add(l)
	}
	if _, ok := u.FishingRod(); ok {
		set.Add(modifier.Fisher)
	}
	if p, ok := u.Phone(); ok {
		set.Add(modifier.Phone)
		if w.SMS.Count(p.Number) > 0 {
			set.Add(modifier.SMS)
		}
	}
	if p, ok := u.Pet(); ok {
		set.Add(petModifier(p))
	}
	return set
}

func luckModifier(l float64) (m *modifier.Mod, ok bool) {
	var x *modifier.Mod
	switch {
	case l < 0.05:
		x = modifier.TerribleLuck
	case l < 0.20:
		x = modifier.BadLuck
	case l > 0.95:
		x = modifier.ExcellentLuck
	case l > 0.80:
		x = modifier.GoodLuck
	default:
		return nil, false
	}
	return &modifier.Mod{
		Emoji: x.Emoji,
		// no need for multiplier, luck is already
		// used in strength calculation
		Multiplier:  0,
		Description: x.Description,
	}, true
}

func petModifier(p *pets.Pet) *modifier.Mod {
	var q float64
	switch p.Species.Quality() {
	case pets.Common:
		q = 0.05
	case pets.Rare:
		q = 0.10
	case pets.Exotic:
		q = 0.15
	case pets.Legendary:
		q = 0.20
	}
	r := ""
	if p.Species.Quality() != pets.Common {
		r = fmt.Sprintf("%s ", p.Species.Quality())
	}
	return &modifier.Mod{
		Emoji:       p.Species.Emoji(),
		Multiplier:  q,
		Description: fmt.Sprintf("У вас есть %sпитомец: <code>%s</code>", r, p),
	}
}
