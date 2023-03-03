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
	"nechego/token"
	"time"
)

// User represents a player.
type User struct {
	TUID        int64        // Telegram ID.
	Energy      Energy       // Energy level.
	Rating      float64      // Elo rating in fights.
	Messages    int          // Number of messages sent.
	BannedUntil time.Time    // Time after which the user is unbanned.
	Status      string       // Status displayed in the profile.
	Inventory   *item.Set    // Personal items.
	Net         *fishing.Net // Net if currently cast, else nil.
	LastMessage time.Time    // When was the last message sent?
	Buffs       buff.Set     // Active buffs.
	Developer   bool         // Flag of a game developer.
	CombatMode  pvp.Mode     // PvP or PvE?
	Funds       Funds        // Collectable items.
	Farm        *farm.Farm   // The source of vegetables.
	Retired     time.Time    // When the job shift should finish.
	Friends     Friends      // The list of friends' TUIDs.
}

func NewUser(tuid int64) *User {
	return &User{
		TUID:      tuid,
		Rating:    1500,
		Inventory: item.NewSet(),
		Buffs:     buff.Set{},
		Funds:     Funds{},
		Farm:      farm.New(2, 3),
		Friends:   Friends{},
	}
}

// ID returns the unique user's ID.
func (u *User) ID() int64 {
	return u.TUID
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
		case *food.Food, *food.Meat:
			return x, u.Eat(x)
		}
	}
	return nil, false
}

// InventorySize returns the user's inventory size.
func (u *User) InventorySize() int {
	n := 20
	if x, ok := u.Inventory.ByType(item.TypeLegacy); ok {
		n += 1 + x.Value.(*token.Legacy).Count
	}
	return n
}

// InventoryFull returns true if the item count in the user's
// inventory exceeds inventory size.
func (u *User) InventoryFull() bool {
	return u.Inventory.Count() > u.InventorySize()
}

// InventoryOverflow returns true if the item count in the user's
// inventory exceeds inventory size by a large margin.
func (u *User) InventoryOverflow() bool {
	const margin = 10
	return u.Inventory.Count() > u.InventorySize()+margin
}

// HasSMS returns true if the user has unread SMS.
func (u *User) HasSMS(w *World) bool {
	p, ok := u.Phone()
	return ok && w.SMS.Count(p.Number) > 0
}

func (u *User) Transfer(to *User, x *item.Item) bool {
	if !x.Transferable || !u.Inventory.Remove(x) {
		return false
	}
	to.Funds.Add("обмен", x)
	return true
}
