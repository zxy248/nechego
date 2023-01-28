package game

import (
	"math/rand"
	"nechego/fishing"
	"nechego/food"
	"nechego/item"
	"time"
)

const (
	// InventorySize is a threshold for applying a debuff.
	InventorySize = 15

	// InventoryCap is a threshold for locking commands related to
	// the inventory system.
	InventoryCap = 20
)

type User struct {
	TUID        int64       // Telegram ID.
	Energy      Energy      // Energy level.
	Rating      float64     // Elo rating in fights.
	Messages    int         // Number of messages sent.
	BannedUntil time.Time   // Time after which the user is unbanned.
	Status      string      // Status displayed in the profile.
	Inventory   *item.Items // Personal items.
}

func NewUser(tuid int64) *User {
	return &User{
		TUID:      tuid,
		Rating:    1500,
		Inventory: item.NewItems(),
	}
}

// Eat restores user's energy and removes the specified item from the
// inventory if it can be eaten. If it cannot be eaten, returns false.
func (u *User) Eat(i *item.Item) bool {
	n, ok := i.Value.(food.Nutritioner)
	if !ok {
		return false
	}
	u.Energy.Add(Energy(n.Nutrition() + 0.04*rand.Float64()))
	u.Inventory.Remove(i)
	return true
}

// EatQuick finds the first sensibly eatable item and calls Eat.
func (u *User) EatQuick() (i *item.Item, ok bool) {
	for _, x := range u.Inventory.List() {
		switch v := x.Value.(type) {
		case *fishing.Fish:
			if v.Price() < 2000 {
				return x, u.Eat(x)
			}
		case *food.Food, *food.Meat:
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

// InventoryFull returns true if the item count is greater than InventorySize.
func (u *User) InventoryFull() bool {
	return u.Inventory.Count() > InventorySize
}

// HasSMS return true if the user has unread SMS.
func (u *User) HasSMS(w *World) bool {
	p, ok := u.Phone()
	return ok && w.SMS.Count(p.Number) > 0
}
