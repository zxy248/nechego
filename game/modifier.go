package game

import (
	"nechego/item"
	"nechego/modifier"
)

func (u *User) Modifiers() modifier.Set {
	set := modifier.Set{}
	ms := []modifier.Moder{Luck(u.Luck()), &u.Energy, u.Balance()}
	for _, b := range u.Buffs.List() {
		ms = append(ms, b)
	}
	switch u.RatingPosition {
	case 0:
		ms = append(ms, modifier.RatingFirst)
	case 1:
		ms = append(ms, modifier.RatingSecond)
	case 2:
		ms = append(ms, modifier.RatingThird)
	}

	// Item modifiers.
	// If the same item type is encountered more than once, the
	// modifier will not be applied.
	seen := map[item.Type]bool{item.TypeFishingRod: false}
	for _, x := range u.Inventory.List() {
		if seen[x.Type] {
			continue
		}
		seen[x.Type] = true
		mer, ok := x.Value.(modifier.Moder)
		if !ok {
			continue
		}
		if m, ok := mer.Mod(); ok {
			set.Add(m)
		}
	}

	// Apply modifiers.
	for _, moder := range ms {
		if m, ok := moder.Mod(); ok {
			set.Add(m)
		}
	}
	return set
}
