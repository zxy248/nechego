package item

import (
	"math/rand"
	"nechego/details"
	"nechego/farm/plant"
	"nechego/money"
	"nechego/token"
)

// Items represents a list of items.
// Entries in the list are accessible by keys.
type Items struct {
	I    []*Item
	keys map[int]*Item
}

// NewItems returns an empty Items.
func NewItems() *Items {
	return &Items{I: []*Item{}}
}

// Add appends the item x to the item list.
func (it *Items) Add(x *Item) {
	it.I = append(it.I, x)
}

// Remove removes the item x from the item list.
func (it *Items) Remove(x *Item) bool {
	for i, v := range it.I {
		if v == x {
			it.I = append(it.I[:i], it.I[i+1:]...)
			return true
		}
	}
	// Item not found.
	return false
}

// Contain returns true if the item x is in the list.
func (it *Items) Contain(x *Item) bool {
	for _, v := range it.I {
		if v == x {
			return true
		}
	}
	return false
}

// ByKey gets an item from the list by the key k.
// If there is no such item, ok is false.
func (it *Items) ByKey(k int) (x *Item, ok bool) {
	x, ok = it.keys[k]
	if !ok || !it.Contain(x) {
		return nil, false
	}
	return x, true
}

// ByType gets the first item of the type t from the list.
// If there is no such item, ok is false.
func (it *Items) ByType(t Type) (x *Item, ok bool) {
	for _, x := range it.List() {
		if x.Type == t {
			return x, true
		}
	}
	return nil, false
}

// updateHotkeys remaps the items.
func (it *Items) updateHotkeys() {
	it.keys = map[int]*Item{}
	for i, v := range it.I {
		it.keys[i] = v
	}
}

// Filter goes through the item list and retains only those for which
// keep is true.
func (it *Items) Filter(keep func(i *Item) bool) {
	n := 0
	for _, x := range it.I {
		if keep(x) {
			it.I[n] = x
			n++
		}
	}
	it.I = it.I[:n]
}

// List returns the filtered item list.
func (it *Items) List() []*Item {
	it.Filter(integral)
	return it.list()
}

// HkList updates the hotkeys and returns the filtered item list.
func (it *Items) HkList() []*Item {
	// Do not use this function internally.
	it.Filter(integral)
	it.updateHotkeys()
	return it.list()
}

// list returns a copy of it.I.
func (it *Items) list() []*Item {
	l := make([]*Item, len(it.I))
	copy(l, it.I)
	return l
}

// Move removes the item from the items it and adds it to the items
// dst. Returns false if the item is not transferable or cannot be
// removed.
func (it *Items) Move(dst *Items, x *Item) bool {
	if !x.Transferable {
		return false
	}
	if !it.Remove(x) {
		return false
	}
	dst.Add(x)
	return true
}

// Trim retains only the last n items of the list.
func (it *Items) Trim(n int) {
	if len(it.I) > n {
		it.I = it.I[len(it.I)-n:]
	}
}

// Random returns a random item from the list.
// If there is no items in the list, ok is false.
func (it *Items) Random() (x *Item, ok bool) {
	items := it.List()
	if len(items) == 0 {
		return nil, false
	}
	return items[rand.Intn(len(items))], true
}

// Count returns the number of items in the list.
func (it *Items) Count() int {
	return len(it.I)
}

// PushFront adds the specified items to the head of the list.
func (it *Items) PushFront(i ...*Item) {
	it.I = append(i, it.I...)
}

// Stack aggregates all mergeable items.
func (it *Items) Stack() {
	removed := map[*Item]bool{}
	items := it.List()
	for i, x := range items {
		if removed[x] {
			continue
		}
		for _, y := range items[i+1:] {
			if removed[y] {
				continue
			}
			if r, ok := Merge(x, y); ok && r != nil {
				removed[r] = true
			}
		}
	}
	for x := range removed {
		it.Remove(x)
	}
}

// Merge attemps to stack two items into a single one.
// On success, the returned item should be removed if not nil.
func Merge(a, b *Item) (remove *Item, ok bool) {
	switch x := a.Value.(type) {
	case *plant.Plant:
		y, ok := b.Value.(*plant.Plant)
		if !ok || x.Type != y.Type {
			return nil, false
		}
		x.Count += y.Count
		return b, true

	case *money.Wallet:
		switch y := b.Value.(type) {
		case *money.Cash:
			x.Money += y.Money
			return b, true
		case *money.Wallet:
			x.Money += y.Money
			return nil, true
		}

	case *money.Cash:
		switch y := b.Value.(type) {
		case *money.Cash:
			x.Money += y.Money
			return b, true
		case *money.Wallet:
			y.Money += x.Money
			return a, true
		}

	case *details.Details:
		y, ok := b.Value.(*details.Details)
		if !ok {
			return nil, false
		}
		x.Count += y.Count
		return b, true

	case *token.Legacy:
		y, ok := b.Value.(*token.Legacy)
		if !ok {
			return nil, false
		}
		x.Count += y.Count
		return b, true
	}
	return nil, false
}

// Split attemps to separate the given item into two.
// On success, the returned item should be added to the inventory.
func Split(i *Item, n int) (j *Item, ok bool) {
	if n <= 0 {
		return nil, false
	}
	switch x := i.Value.(type) {
	case *plant.Plant:
		if x.Count <= n {
			return nil, false
		}
		x.Count -= n
		return New(&plant.Plant{Type: x.Type, Count: n}), true

	case *money.Wallet:
		// Not `<=`: allow empty wallet.
		if x.Money < n {
			return nil, false
		}
		x.Money -= n
		return New(&money.Cash{Money: n}), true

	case *money.Cash:
		if x.Money <= n {
			return nil, false
		}
		x.Money -= n
		return New(&money.Cash{Money: n}), true

	case *details.Details:
		if x.Count <= n {
			return nil, false
		}
		x.Count -= n
		return New(&details.Details{Count: n}), true

	case *token.Legacy:
		if x.Count <= n {
			return nil, false
		}
		x.Count -= n
		return New(&token.Legacy{Count: n}), true
	}
	return nil, false
}
