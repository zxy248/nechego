package game

import (
	"math/rand"
	"nechego/buff"
	"nechego/farm"
	"nechego/farm/plant"
	"nechego/fishing"
	"nechego/food"
	"nechego/game/pvp"
	"nechego/item"
	"time"
)

// The number of items in the user's inventory for applying various
// types of restrictions.
const (
	InventorySize = 15
	InventoryCap  = 20
)

// User represents a player.
type User struct {
	TUID        int64        // Telegram ID.
	Energy      Energy       // Energy level.
	Rating      float64      // Elo rating in fights.
	Messages    int          // Number of messages sent.
	BannedUntil time.Time    // Time after which the user is unbanned.
	Status      string       // Status displayed in the profile.
	Inventory   *item.Items  // Personal items.
	Net         *fishing.Net // Net if currently cast, else nil.
	LastMessage time.Time    // When was the last message sent?
	Buffs       buff.Set     // Active buffs.
	Developer   bool         // Flag of a game developer.
	CombatMode  pvp.Mode     // PvP or PvE?
	Funds       Funds        // Collectable items.
	Farm        *farm.Farm   // The source of vegetables.
}

func NewUser(tuid int64) *User {
	return &User{
		TUID:      tuid,
		Rating:    1500,
		Inventory: item.NewItems(),
		Buffs:     buff.Set{},
		Funds:     Funds{},
		Farm:      farm.New(2, 3),
	}
}

// Eat restores user's energy and removes the specified item from the
// inventory if it can be eaten. If it cannot be eaten, returns false.
func (u *User) Eat(i *item.Item) bool {
	n, ok := i.Value.(food.Nutritioner)
	if !ok {
		return false
	}
	switch x := i.Value.(type) {
	case *plant.Plant:
		if x.Count > 1 {
			x.Count--
		} else {
			u.Inventory.Remove(i)
		}
	default:
		u.Inventory.Remove(i)
	}
	u.Energy.Add(Energy(n.Nutrition() + 0.04*rand.Float64()))
	if f, ok := i.Value.(*food.Food); ok {
		if f.Type == food.Beer {
			u.Buffs.Apply(buff.Beer, 10*time.Minute)
		}
	}
	return true
}

// EatQuick finds the first sensibly eatable item and calls Eat.
// If there is no food in the inventory, returns (nil, false).
func (u *User) EatQuick() (i *item.Item, ok bool) {
	for _, x := range u.Inventory.List() {
		switch v := x.Value.(type) {
		case *fishing.Fish:
			if v.Price() < 2000 {
				return x, u.Eat(x)
			}
		case *food.Food, *food.Meat, *plant.Plant:
			return x, u.Eat(x)
		}
	}
	return nil, false
}

// InventoryFull returns true if the item count in the user's
// inventory exceeds inventory size.
func (u *User) InventoryFull() bool {
	return u.Inventory.Count() > InventorySize
}

// HasSMS returns true if the user has unread SMS.
func (u *User) HasSMS(w *World) bool {
	p, ok := u.Phone()
	return ok && w.SMS.Count(p.Number) > 0
}
