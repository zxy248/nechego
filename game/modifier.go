package game

import (
	"nechego/item"
	"nechego/modifier"
)

func (u *User) Modset(w *World) modifier.Set {
	set := modifier.Set{}
	moders := []modifier.Moder{
		Luck(u.Luck()),
		&u.Energy,
		u.Balance(),
	}

	// If the predicate is true, the corresponding modifier will
	// be added to the set.
	table := []struct {
		predicate func() bool
		modifier  *modifier.Mod
	}{
		{u.InventoryFull, modifier.Heavy},
		{func() bool { return u.HasSMS(w) }, modifier.SMS},
	}
	for _, x := range table {
		if x.predicate() {
			moders = append(moders, x.modifier)
		}
	}

	// Item modifiers.
	// If the same item type is encountered more than once, the
	// modifier will not be applied.
	seen := map[item.Type]bool{}
	for _, x := range u.Inventory.List() {
		if seen[x.Type] {
			continue
		}
		seen[x.Type] = true

		moder, ok := x.Value.(modifier.Moder)
		if !ok {
			continue
		}
		if m, ok := moder.Mod(); ok {
			set.Add(m)
		}
	}

	// Apply modifiers.
	for _, moder := range moders {
		if m, ok := moder.Mod(); ok {
			set.Add(m)
		}
	}
	return set
}
