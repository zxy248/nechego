package game

import (
	"nechego/farm/plant"
	"nechego/item"
)

// Plant plants the specified item at the Farm if possible.
func (u *User) Plant(i *item.Item) *plant.Plant {
	p, ok := i.Value.(*plant.Plant)
	if !ok {
		return &plant.Plant{}
	}

	var c int
	for c = 0; c < p.Count; c++ {
		if !u.Farm.Plant(p.Type) {
			break
		}
	}
	p.Count -= c
	if p.Count == 0 && !u.Inventory.Remove(i) {
		panic("cannot remove zero Plant from inventory")
	}
	return &plant.Plant{Type: p.Type, Count: c}
}

// Harvest harvests all grown plants and adds them to the inventory.
// Returns harvested plants.
func (u *User) Harvest() []*plant.Plant {
	harvested := u.Farm.Harvest()
	for _, p := range harvested {
		u.Inventory.Add(item.New(p))
	}
	return harvested
}

func (u *User) FarmUpgradeCost() int {
	return 1000 * fib(u.Farm.Columns+u.Farm.Rows)
}

func (u *User) UpgradeFarm() bool {
	if !u.Balance().Spend(u.FarmUpgradeCost()) {
		return false
	}
	u.Farm.Grow()
	return true
}

var fibCache = map[int]int{0: 0, 1: 1}

func fib(n int) int {
	if n < 0 {
		panic("fib negative argument")
	}
	if k, ok := fibCache[n]; ok {
		return k
	}
	k := fib(n-1) + fib(n-2)
	fibCache[n] = k
	return k
}
