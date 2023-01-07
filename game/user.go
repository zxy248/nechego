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
	TUID      int64
	Energy    int
	EnergyCap int
	Rating    float64
	Messages  int
	Banned    bool
	Birthday  time.Time
	Gender    Gender
	Status    string
	Inventory []*Item

	hotkeys map[string]int // map from hotkey to item ID
}

func NewUser(tuid int64) *User {
	return &User{
		TUID:      tuid,
		EnergyCap: 5,
		Rating:    1500,
		Inventory: []*Item{},

		hotkeys: map[string]int{},
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

func (u *User) GenerateHotkeys() {
	u.hotkeys = map[string]int{}
	for i, v := range u.Inventory {
		u.hotkeys[string(Hotkeys[i])] = v.ID
	}
}

func (u *User) ItemByID(id int) (i *Item, ok bool) {
	for _, v := range u.Inventory {
		if v.ID == id {
			return v, true
		}
	}
	return nil, false
}

func (u *User) ItemByHotkey(h string) (i *Item, ok bool) {
	return u.ItemByID(u.hotkeys[h])
}

func (u *User) TraverseInventory(f func(*Item)) {
	for _, v := range u.Inventory {
		f(v)
	}
}

func (u *User) Total() int {
	t := 0
	u.TraverseInventory(func(i *Item) {
		switch o := i.Object.(type) {
		case *Wallet:
			t += o.Money
		case *CreditCard:
			t += o.Money
		case *Debt:
			t -= o.Money
		}
	})
	return t
}

func (u *User) InDebt() bool {
	f := false
	u.TraverseInventory(func(i *Item) {
		switch i.Object.(type) {
		case Debt:
			f = true
			return
		}
	})
	return f
}

func (u *User) IsEblan() bool {
	f := false
	u.TraverseInventory(func(i *Item) {
		switch i.Object.(type) {
		case EblanToken:
			f = true
			return
		}
	})
	return f
}

func (u *User) IsAdmin() bool {
	f := false
	u.TraverseInventory(func(i *Item) {
		switch i.Object.(type) {
		case AdminToken:
			f = true
			return
		}
	})
	return f
}

func (u *User) IsPair() bool {
	f := false
	u.TraverseInventory(func(i *Item) {
		switch i.Object.(type) {
		case PairToken:
			f = true
			return
		}
	})
	return f
}
