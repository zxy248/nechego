package game

import (
	"errors"
	"math/rand"
	"nechego/fishing"
	"nechego/item"
)

// FishCatchProb returns the user's chance to catch fish.
func (u *User) FishCatchProb() float64 {
	p := 0.5
	p += -0.04 + 0.08*u.Luck()
	return p
}

// Fish returns a new random item to be added to the user's inventory
// and decreases durability of the fishing rod r. If the user was not
// able to catch any fish, returns (nil, false).
func (u *User) Fish(r *fishing.Rod) (i *item.Item, ok bool) {
	r.Durability -= 0.01
	if rand.Float64() > u.FishCatchProb() {
		return nil, false
	}
	if rand.Float64() < 0.08 {
		return item.Random(), true
	}
	f := fishing.RandomFish()
	quality := 1 + 0.1*float64(r.Level)
	luck := 0.9 + 0.2*u.Luck()
	multiplier := quality * luck
	f.Length *= multiplier
	f.Weight *= multiplier
	return item.New(f), true
}

var (
	ErrNoNet          = errors.New("no fishing net in inventory")
	ErrNetAlreadyCast = errors.New("fishing net is already cast")
)

// CastNet removes the fishing net from the user's inventory and casts it.
func (u *User) CastNet() error {
	if u.Net != nil {
		return ErrNetAlreadyCast
	}
	var net *fishing.Net
	var ok bool
	u.Inventory.Pop(func(x *item.Item) bool {
		net, ok = x.Value.(*fishing.Net)
		return ok
	})
	if !ok {
		return ErrNoNet
	}
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
