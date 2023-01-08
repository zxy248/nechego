package game

import (
	"time"
)

type Gender int

const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
	GenderTrans
)

type User struct {
	TUID             int64
	Energy           int
	EnergyCap        int
	Rating           float64
	Messages         int
	Banned           bool
	Birthday         time.Time
	Gender           Gender
	Status           string
	Inventory        []*Item
	inventoryHotkeys map[int]*Item
}

func NewUser(tuid int64) *User {
	return &User{
		TUID:      tuid,
		EnergyCap: 5,
		Rating:    1500,
		Inventory: []*Item{},
	}
}

func (u *User) Ban() {
	u.Banned = true
}

func (u *User) Unban() {
	u.Banned = false
}

func (u *User) IncrementMessages() {
	u.Messages++
}

func (u *User) AddRating(δ float64) {
	u.Rating += δ
}

func (u *User) SpendEnergy(δ int) bool {
	if u.Energy < δ {
		return false
	}
	u.Energy -= δ
	return true
}

func (u *User) RestoreEnergy(δ int) {
	u.Energy += δ
	if u.Energy > u.EnergyCap {
		u.Energy = u.EnergyCap
	}
}

func (u *User) Items() []*Item {
	n := 0
	for _, v := range u.Inventory {
		if v.Expire.IsZero() || time.Now().Before(v.Expire) {
			u.Inventory[n] = v
			n++
		}
	}
	u.Inventory = u.Inventory[:n]
	return u.Inventory
}

func (u *User) ListInventory() []*Item {
	var r []*Item
	u.inventoryHotkeys, r = hotkeys(u.Inventory)
	return r
}

func (u *User) ItemByID(id int) (i *Item, ok bool) {
	for _, v := range u.Inventory {
		if v.ID == id {
			return v, true
		}
	}
	return nil, false
}

func (u *User) HasItem(i *Item) bool {
	for _, j := range u.Items() {
		if i == j {
			return true
		}
	}
	return false
}

func (u *User) ItemByHotkey(k int) (i *Item, ok bool) {
	i, ok = u.inventoryHotkeys[k]
	if !ok || !u.HasItem(i) {
		return nil, false
	}
	return i, true
}

func (u *User) Total() int {
	t := 0
	for _, v := range u.Items() {
		switch x := v.Value.(type) {
		case *Wallet:
			t += x.Money
		case *CreditCard:
			t += x.Money
		case *Debt:
			t -= x.Money
		}
	}
	return t
}

func (u *User) InDebt() bool {
	for _, v := range u.Items() {
		switch v.Value.(type) {
		case *Debt:
			return true
		}
	}
	return false
}

func (u *User) IsEblan() bool {
	for _, v := range u.Items() {
		switch v.Value.(type) {
		case *EblanToken:
			return true
		}
	}
	return false
}

func (u *User) IsAdmin() bool {
	for _, v := range u.Items() {
		switch v.Value.(type) {
		case *AdminToken:
			return true
		}
	}
	return false
}

func (u *User) IsPair() bool {
	for _, v := range u.Items() {
		switch v.Value.(type) {
		case *PairToken:
			return true
		}
	}
	return false
}
