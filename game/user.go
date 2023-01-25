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
	InventorySize = 15
	InventoryCap  = 20
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
	case *food.Meat:
		u.Inventory.Remove(i)
		switch x.Species.Size() {
		case pets.Small:
			u.RestoreEnergy(10 + rand.Intn(10))
		case pets.Medium:
			u.RestoreEnergy(20 + rand.Intn(10))
		case pets.Big:
			u.RestoreEnergy(30 + rand.Intn(10))
		}
		return true
	}
	return false
}

func (u *User) EatQuick() (i *item.Item, ok bool) {
	for _, x := range u.Inventory.List() {
		switch v := x.Value.(type) {
		case *fishing.Fish:
			if v.Price() < 2000 {
				return x, u.Eat(x)
			}
		case *food.Food:
			return x, u.Eat(x)
		case *food.Meat:
			return x, u.Eat(x)
		}
	}
	return nil, false
}

func (u *User) Fish(r *fishing.Rod) *item.Item {
	r.Durability -= 0.01
	if rand.Float64() < 0.08 {
		return item.Random()
	}

	quality := 1.0 + 0.5*float64(r.Quality)
	luck := 0.9 + 0.2*u.Luck()
	total := quality * luck

	f := fishing.RandomFish()
	f.Length *= total
	f.Weight *= total
	return &item.Item{Type: item.TypeFish, Transferable: true, Value: f}
}
