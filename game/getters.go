package game

import (
	"nechego/money"
	"nechego/pets"
	"nechego/token"
)

func (u *User) Cash() (c *money.Cash, ok bool) {
	for _, x := range u.Inventory.Normal() {
		if c, ok := x.Value.(*money.Cash); ok {
			return c, true
		}
	}
	return nil, false
}

func (u *User) Wallet() (w *money.Wallet, ok bool) {
	for _, x := range u.Inventory.Normal() {
		if w, ok := x.Value.(*money.Wallet); ok {
			return w, true
		}
	}
	return nil, false
}

func (u *User) Dice() (d *token.Dice, ok bool) {
	for _, x := range u.Inventory.Normal() {
		if d, ok = x.Value.(*token.Dice); ok {
			return d, true
		}
	}
	return nil, false
}

func (u *User) Eblan() bool {
	for _, x := range u.Inventory.Normal() {
		if _, ok := x.Value.(*token.Eblan); ok {
			return true
		}
	}
	return false
}

func (u *User) Admin() bool {
	for _, x := range u.Inventory.Normal() {
		if _, ok := x.Value.(*token.Admin); ok {
			return true
		}
	}
	return false
}

func (u *User) Pair() bool {
	for _, x := range u.Inventory.Normal() {
		if _, ok := x.Value.(*token.Pair); ok {
			return true
		}
	}
	return false
}

func (u *User) Pet() (p *pets.Pet, ok bool) {
	for _, x := range u.Inventory.Normal() {
		if p, ok = x.Value.(*pets.Pet); ok {
			return p, true
		}
	}
	return nil, false
}
