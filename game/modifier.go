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
	}
	for _, x := range table {
		if x.predicate() {
			moders = append(moders, x.modifier)
		}
	}

	// Buff modifiers.
	for _, b := range u.Buffs.List() {
		moders = append(moders, b)
	}

	// Rating modifiers.
	if top := w.SortedUsers(ByElo); len(top) >= 3 {
		switch u {
		case top[0]:
			moders = append(moders, modifier.RatingFirst)
		case top[1]:
			moders = append(moders, modifier.RatingSecond)
		case top[2]:
			moders = append(moders, modifier.RatingThird)
		}
	}

	// Item modifiers.
	// If the same item type is encountered more than once, the
	// modifier will not be applied.
	seen := map[item.Type]int{item.TypeFishingRod: 0}
	for _, x := range u.Inventory.List() {
		if n, ok := seen[x.Type]; ok {
			if n > 0 {
				continue
			}
			seen[x.Type]++
		}
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
