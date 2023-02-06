package game

import (
	"errors"
	"math/rand"
	"nechego/fishing"
	"nechego/food"
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

// Fish returns a new random item to be added to the user's inventory
// and decreases durability of the fishing rod r.
func (u *User) Fish(r *fishing.Rod) *item.Item {
	r.Durability -= 0.01
	if rand.Float64() < 0.08 {
		return item.Random()
	}

	quality := 1 + 0.1*float64(r.Level)
	luck := 0.9 + 0.2*u.Luck()
	total := quality * luck

	f := fishing.RandomFish()
	f.Length *= total
	f.Weight *= total
	return item.New(f)
}

var (
	ErrNoNet          = errors.New("no fishing net in inventory")
	ErrNetAlreadyCast = errors.New("fishing net is already cast")
	ErrFishInNet      = errors.New("there is fish in fishing net")
)

// CastNet removes the fishing net from the user's inventory and casts it.
func (u *User) CastNet() error {
	if u.Net != nil {
		return ErrNetAlreadyCast
	}
	x, ok := u.Inventory.ByType(item.TypeFishingNet)
	if !ok {
		return ErrNoNet
	}
	net := x.Value.(*fishing.Net)
	if net.Count() != 0 {
		return ErrFishInNet
	}
	u.Inventory.Remove(x)
	u.Net = net
	return nil
}

// DrawNet returns the fishing net to the user's inventory if it is
// currently cast.
func (u *User) DrawNew() (n *fishing.Net, ok bool) {
	if u.Net == nil {
		return nil, false
	}
	net := u.Net
	u.Net = nil
	u.Inventory.Add(item.New(net))
	return net, true
}

// FillNet adds a random fish to the cast fishing net.
func (u *User) FillNet() {
	if u.Net == nil {
		return
	}
	u.Net.Fill()
}

// UnloadNet moves the caught fish from the specified fishing net to
// the user's inventory.
func (u *User) UnloadNet(n *fishing.Net) []*fishing.Fish {
	catch := n.Unload()
	for _, f := range catch {
		u.Inventory.Add(item.New(f))
	}
	return catch
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
