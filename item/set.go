package item

import (
	"math/rand"
	"nechego/details"
	"nechego/farm/plant"
	"nechego/money"
	"nechego/token"
)

// Set represents a bunch of items.
// Entries in the set are indexed by integer keys.
type Set struct {
	I    []*Item
	keys map[int]*Item
}

// NewSet returns an empty Set.
func NewSet() *Set {
	return &Set{I: []*Item{}}
}

// Add adds the specified items to the tail of the set.
func (s *Set) Add(x ...*Item) {
	s.I = append(s.I, x...)
}

// AddFront adds the specified items to the head of the set.
func (s *Set) AddFront(x ...*Item) {
	s.I = append(x, s.I...)
}

// Remove removes the item x from the set.
func (s *Set) Remove(x *Item) bool {
	for i, v := range s.I {
		if v == x {
			s.I = append(s.I[:i], s.I[i+1:]...)
			return true
		}
	}
	// Item not found.
	return false
}

// Contain returns true if the item x is in the set.
func (s *Set) Contain(x *Item) bool {
	for _, v := range s.List() {
		if v == x {
			return true
		}
	}
	return false
}

// ByKey gets an item from the set by the key k.
// If there is no such item, returns (nil, false).
func (s *Set) ByKey(k int) (x *Item, ok bool) {
	x, ok = s.keys[k]
	if !ok || !s.Contain(x) {
		return nil, false
	}
	return x, true
}

// ByType gets the first item of the type t from the set.
// If there is no such item, returns (nil, false).
func (s *Set) ByType(t Type) (x *Item, ok bool) {
	for _, x := range s.List() {
		if x.Type == t {
			return x, true
		}
	}
	return nil, false
}

// indexItems indexes all items in the set.
func (s *Set) indexItems() {
	s.keys = map[int]*Item{}
	for i, v := range s.I {
		s.keys[i] = v
	}
}

// Filter retains only those items in the set for which keep is true.
func (s *Set) Filter(keep func(i *Item) bool) {
	n := 0
	for _, x := range s.I {
		if keep(x) {
			s.I[n] = x
			n++
		}
	}
	s.I = s.I[:n]
}

// List returns the filtered list of all items in the set.
func (s *Set) List() []*Item {
	s.Filter(integral)
	return s.itemsCopy()
}

// HkList returns the filtered, reindexed list of all items in the set.
func (s *Set) HkList() []*Item {
	// Do not use this function internally.
	s.Filter(integral)
	s.indexItems()
	return s.itemsCopy()
}

// itemsCopy returns a copy of the set's item list.
func (s *Set) itemsCopy() []*Item {
	l := make([]*Item, len(s.I))
	copy(l, s.I)
	return l
}

// Move removes the item from the set s and adds it to the set to.
// Returns false if the item is not transferable or cannot be removed.
func (s *Set) Move(to *Set, x *Item) bool {
	if !x.Transferable || !s.Remove(x) {
		return false
	}
	to.Add(x)
	return true
}

// Trim retains the last n items in the set.
func (s *Set) Trim(n int) {
	if len(s.I) > n {
		s.I = s.I[len(s.I)-n:]
	}
}

// Random returns a random item from the set.
// If there is no items in the set, returns (nil, false).
func (s *Set) Random() (x *Item, ok bool) {
	items := s.List()
	if len(items) == 0 {
		return nil, false
	}
	return items[rand.Intn(len(items))], true
}

// Count returns the number of items in the set.
func (s *Set) Count() int {
	return len(s.I)
}

// Stack aggregates all mergeable items.
func (s *Set) Stack() {
	removed := map[*Item]bool{}
	items := s.List()
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
		s.Remove(x)
	}
}

// Merge attemps to stack two items into a single slot.
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

// Split attemps to separate the given item into two items.
// On success, the returned item should be added to the inventory.
func Split(a *Item, n int) (b *Item, ok bool) {
	if n <= 0 {
		return nil, false
	}
	switch x := a.Value.(type) {
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
