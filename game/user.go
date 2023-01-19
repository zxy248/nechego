package game

import (
	"math/rand"
	"nechego/fishing"
	"nechego/food"
	"nechego/item"
	"nechego/pets"
	"time"
)

const (
	InventorySize = 10
	InventoryCap  = 17
)

type User struct {
	TUID      int64
	Energy    int
	Rating    float64
	Messages  int
	Banned    bool
	Birthday  time.Time
	Gender    Gender
	Status    string
	Inventory *item.Items
}

func NewUser(tuid int64) *User {
	return &User{
		TUID:      tuid,
		Rating:    1500,
		Inventory: item.NewItems(),
	}
}

func (u *User) Eat(i *item.Item) bool {
	switch x := i.Value.(type) {
	case *fishing.Fish:
		u.Inventory.Remove(i)
		if x.Heavy() {
			u.RestoreEnergy(25 + rand.Intn(10))
		} else {
			u.RestoreEnergy(15 + rand.Intn(10))
		}
		return true
	case *pets.Pet:
		if x.Name != "" {
			return false
		}
		u.Inventory.Remove(i)
		switch x.Species.Size() {
		case pets.Small:
			u.RestoreEnergy(5 + rand.Intn(10))
		case pets.Medium:
			u.RestoreEnergy(15 + rand.Intn(10))
		case pets.Big:
			u.RestoreEnergy(25 + rand.Intn(10))
		}
		return true
	case *food.Food:
		u.Inventory.Remove(i)
		u.RestoreEnergy(int(x.Nutrition() * 100))
		return true
	}
	return false
}

func (u *User) EatQuick() (i *item.Item, ok bool) {
	for _, i = range u.Inventory.Normal() {
		switch x := i.Value.(type) {
		case *fishing.Fish:
			if x.Price() < 2000 {
				return i, u.Eat(i)
			}
		case *food.Food:
			return i, u.Eat(i)
		}
	}
	return nil, false
}
