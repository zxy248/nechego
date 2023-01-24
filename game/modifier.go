package game

import (
	"fmt"
	"nechego/modifier"
	"nechego/pets"
)

func (u *User) Modset(w *World) modifier.Set {
	table := []struct {
		predicate func() bool
		modifier  *modifier.Mod
	}{
		{u.Admin, modifier.Admin},
		{u.Eblan, modifier.Eblan},
		{u.Rich, modifier.Rich},
		{u.Poor, modifier.Poor},
		{func() bool { return u.Energy < EnergyCap/10 }, modifier.NoEnergy},
		{func() bool { return u.Energy == EnergyCap }, modifier.FullEnergy},
		{func() bool { return u.Inventory.Count() > InventorySize }, modifier.Heavy},
		{func() bool { _, ok := u.FishingRod(); return ok }, modifier.Fisher},
		{func() bool { _, ok := u.Phone(); return ok }, modifier.Phone},
		{func() bool { p, ok := u.Phone(); return ok && w.SMS.Count(p.Number) > 0 }, modifier.SMS},
	}
	set := modifier.Set{}
	for _, x := range table {
		if x.predicate() {
			set.Add(x.modifier)
		}
	}
	if l, ok := luckModifier(u.Luck()); ok {
		set.Add(l)
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
		Emoji:       x.Emoji,
		Description: x.Description,
		// no need for multiplier, luck is already
		// used in strength calculation
	}, true
}

func petModifier(p *pets.Pet) *modifier.Mod {
	var multiplier float64
	switch p.Species.Quality() {
	case pets.Common:
		multiplier = 0.05
	case pets.Rare:
		multiplier = 0.10
	case pets.Exotic:
		multiplier = 0.15
	case pets.Legendary:
		multiplier = 0.20
	}
	prefix := ""
	if p.Species.Quality() != pets.Common {
		prefix = fmt.Sprintf("%s ", p.Species.Quality())
	}
	return &modifier.Mod{
		Emoji:       p.Species.Emoji(),
		Multiplier:  multiplier,
		Description: fmt.Sprintf("У вас есть %sпитомец: <code>%s</code>", prefix, p),
	}
}
