package game

import (
	"nechego/fishing"
	"nechego/money"
	"nechego/pets"
	"nechego/phone"
	"nechego/token"
)

func (u *User) Cash() (c *money.Cash, ok bool) {
	for _, x := range u.Inventory.List() {
		if c, ok := x.Value.(*money.Cash); ok {
			return c, true
		}
	}
	return nil, false
}

func (u *User) Wallet() (w *money.Wallet, ok bool) {
	for _, x := range u.Inventory.List() {
		if w, ok := x.Value.(*money.Wallet); ok {
			return w, true
		}
	}
	return nil, false
}

func (u *User) Dice() (d *token.Dice, ok bool) {
	for _, x := range u.Inventory.List() {
		if d, ok = x.Value.(*token.Dice); ok {
			return d, true
		}
	}
	return nil, false
}

func (u *User) Eblan() bool {
	for _, x := range u.Inventory.List() {
		if _, ok := x.Value.(*token.Eblan); ok {
			return true
		}
	}
	return false
}

func (u *User) Admin() bool {
	for _, x := range u.Inventory.List() {
		if _, ok := x.Value.(*token.Admin); ok {
			return true
		}
	}
	return false
}

func (u *User) Pair() bool {
	for _, x := range u.Inventory.List() {
		if _, ok := x.Value.(*token.Pair); ok {
			return true
		}
	}
	return false
}

func (u *User) Pet() (p *pets.Pet, ok bool) {
	for _, x := range u.Inventory.List() {
		if p, ok = x.Value.(*pets.Pet); ok {
			return p, true
		}
	}
	return nil, false
}

func (u *User) FishingRod() (r *fishing.Rod, ok bool) {
	for _, x := range u.Inventory.List() {
		if r, ok := x.Value.(*fishing.Rod); ok {
			return r, true
		}
	}
	return nil, false
}

func (u *User) Phone() (p *phone.Phone, ok bool) {
	for _, x := range u.Inventory.List() {
		if p, ok := x.Value.(*phone.Phone); ok {
			return p, true
		}
	}
	return nil, false
}
