package game

import (
	"nechego/buff"
	"nechego/farm/plant"
	"nechego/fishing"
	"nechego/food"
	"nechego/item"
	"nechego/pets"
	"time"
)

// EatQuick finds the first sensibly eatable item and calls Eat.
// If there is no food in the inventory, returns (nil, false).
func (u *User) EatQuick() (i *item.Item, ok bool) {
	for _, x := range u.Inventory.List() {
		if quickEatable(x) {
			return x, u.Eat(x)
		}
	}
	return nil, false
}

func quickEatable(x *item.Item) bool {
	switch v := x.Value.(type) {
	case *food.Food, *food.Meat:
		return true
	case *fishing.Fish:
		return v.Price() < 2000
	default:
		return false
	}
}

// Eat restores user's energy and removes the specified item from the
// inventory if it can be eaten. If it cannot be eaten, returns false.
func (u *User) Eat(x *item.Item) bool {
	var keep bool
	switch x := x.Value.(type) {
	case *plant.Plant:
		x.Count--
		keep = x.Count > 0
		eatPlant(u, x)
	case *fishing.Fish:
		eatFish(u, x)
	case *food.Food:
		eatFood(u, x)
	case *food.Meat:
		eatMeat(u, x)
	case *pets.Pet:
		eatPet(u, x)
	default:
		return false
	}
	if !keep {
		u.Inventory.Remove(x)
	}
	return true
}

func eatPlant(u *User, p *plant.Plant) {
	u.Energy.Add(0.1)
}

func eatFish(u *User, f *fishing.Fish) {
	if f.Heavy() {
		u.Energy.Add(0.25)
	} else {
		u.Energy.Add(0.15)
	}
}

func eatFood(u *User, f *food.Food) {
	u.Energy.Add(Energy(f.Nutrition()))
	if f.Type == food.Beer {
		u.Buffs.Apply(buff.Beer, 10*time.Minute)
	}
}

func eatMeat(u *User, m *food.Meat) {
	switch m.Species.Size() {
	case pets.Small:
		u.Energy.Add(0.1)
	case pets.Medium:
		u.Energy.Add(0.2)
	case pets.Big:
		u.Energy.Add(0.3)
	}
}

func eatPet(u *User, p *pets.Pet) {
	switch p.Species.Size() {
	case pets.Small:
		u.Energy.Add(0.05)
	case pets.Medium:
		u.Energy.Add(0.15)
	case pets.Big:
		u.Energy.Add(0.25)
	}
}
