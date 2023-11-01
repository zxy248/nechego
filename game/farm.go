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

const (
	MaxFarmRows    = 8
	MaxFarmColumns = 8
)

// FarmUpgradeCost returns the cost of the farm expansion.
// If the farm cannot be upgraded, returns (0, false).
func (u *User) FarmUpgradeCost() (cost int, ok bool) {
	if u.Farm.Rows >= MaxFarmRows && u.Farm.Columns >= MaxFarmColumns {
		return 0, false
	}
	return 1000 * fibonacci(u.Farm.Columns+u.Farm.Rows), true
}

func (u *User) UpgradeFarm(cost int) bool {
	if !u.Balance().Spend(cost) {
		return false
	}
	u.Farm.Grow()
	return true
}

var fibonacciCache = map[int]int{0: 0, 1: 1}

func fibonacci(n int) int {
	if n < 0 {
		panic("fibonacci negative argument")
	}
	if k, ok := fibonacciCache[n]; ok {
		return k
	}
	k := fibonacci(n-1) + fibonacci(n-2)
	fibonacciCache[n] = k
	return k
}
