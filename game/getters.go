package game

import (
	"nechego/fishing"
	"nechego/pets"
	"nechego/token"
)

// Spender is implemented by any value that can be partially spended.
type Spender interface {
	// Spend must return true on success, or false if the value
	// cannot be spended by n.
	Spend(n int) bool
}

func GetItem[T any](u *User) (x T, ok bool) {
	for _, x := range u.Inventory.List() {
		if x, ok := x.Value.(T); ok {
			return x, true
		}
	}
	return x, false
}

func (u *User) Eblan() bool { _, ok := GetItem[*token.Eblan](u); return ok }
func (u *User) Admin() bool { _, ok := GetItem[*token.Admin](u); return ok }
func (u *User) Pair() bool  { _, ok := GetItem[*token.Pair](u); return ok }
func (u *User) Dice() bool  { _, ok := GetItem[*token.Dice](u); return ok }

func (u *User) Pet() (p *pets.Pet, ok bool)           { return GetItem[*pets.Pet](u) }
func (u *User) FishingRod() (r *fishing.Rod, ok bool) { return GetItem[*fishing.Rod](u) }
func (u *User) FishingNet() (n *fishing.Net, ok bool) { return GetItem[*fishing.Net](u) }
