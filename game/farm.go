package game

import (
	"nechego/farm"
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

// Pick pops the Plant at the specified location and adds it to the
// user's inventory.
func (u *User) PickPlant(row, column int) (p *plant.Plant, ok bool) {
	p, ok = u.Farm.Pick(farm.Plot{Row: row, Column: column})
	if !ok {
		return nil, false
	}
	u.Inventory.Add(item.New(p))
	return p, true
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
	return 1000 * fib(u.Farm.Columns+u.Farm.Rows), true
}

func (u *User) UpgradeFarm(cost int) bool {
	if !u.Balance().Spend(cost) {
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
