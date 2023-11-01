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

// Remove removes the specified item from the set and returns true.
// If the item is not in the set, returns false.
func (s *Set) Remove(x *Item) bool {
	_, ok := s.Pop(func(y *Item) bool {
		return x == y
	})
	return ok
}

// Pop removes and returns the first item from the list for which the
// specified predicate p is true. If there is no such item, returns
// (nil, false).
func (s *Set) Pop(p func(*Item) bool) (x *Item, ok bool) {
	for i, v := range s.I {
		if p(v) {
			s.I = append(s.I[:i], s.I[i+1:]...)
			return v, true
		}
	}
	return nil, false
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

// indexItems indexes all items in the set.
func (s *Set) indexItems() {
	s.keys = map[int]*Item{}
	for i, v := range s.I {
		s.keys[i] = v
	}
}

// Keep retains only those items in the set for which the specified
// predicate p is true.
func (s *Set) Keep(p func(i *Item) bool) {
	n := 0
	for _, x := range s.I {
		if p(x) {
			s.I[n] = x
			n++
		}
	}
	s.I = s.I[:n]
}

// List returns the filtered list of all items in the set.
func (s *Set) List() []*Item {
	s.Keep(integral)
	return s.itemsCopy()
}

// HkList returns the filtered, reindexed list of all items in the set.
func (s *Set) HkList() []*Item {
	// Do not use this function internally.
	s.Keep(integral)
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
	for i, a := range items {
		if removed[a] {
			continue
		}
		for _, b := range items[i+1:] {
			if removed[b] {
				continue
			}
			if Merge(a, b) {
				removed[b] = true
			}
		}
	}
	for x := range removed {
		s.Remove(x)
	}
}

// Merge attemps to merge the item b into the item a.
// On success, the item b must be removed from the set.
func Merge(a, b *Item) bool {
	switch x := a.Value.(type) {
	case *plant.Plant:
		return mergePlant(x, b.Value)
	case *money.Wallet:
		return mergeWallet(x, b.Value)
	case *money.Cash:
		return mergeCash(x, b.Value)
	case *details.Details:
		return mergeDetails(x, b.Value)
	case *token.Legacy:
		return mergeLegacy(x, b.Value)
	}
	return false
}

func mergePlant(p *plant.Plant, v any) bool {
	if x, ok := v.(*plant.Plant); ok && p.Type == x.Type {
		p.Count += x.Count
		return true
	}
	return false
}

func mergeWallet(w *money.Wallet, v any) bool {
	if x, ok := v.(*money.Cash); ok {
		w.Money += x.Money
		return true
	}
	return false
}

func mergeCash(c *money.Cash, v any) bool {
	if x, ok := v.(*money.Cash); ok {
		c.Money += x.Money
		return true
	}
	return false
}

func mergeDetails(d *details.Details, v any) bool {
	if x, ok := v.(*details.Details); ok {
		d.Count += x.Count
		return true
	}
	return false
}

func mergeLegacy(l *token.Legacy, v any) bool {
	if x, ok := v.(*token.Legacy); ok {
		l.Count += x.Count
		return true
	}
	return false
}

// Split attemps to separate the given item into two items.
// On success, the returned item should be added to the inventory.
func Split(a *Item, n int) (b *Item, ok bool) {
	if n <= 0 {
		return nil, false
	}
	switch x := a.Value.(type) {
	case *plant.Plant:
		return splitPlant(x, n)
	case *money.Wallet:
		return splitWallet(x, n)
	case *money.Cash:
		return splitCash(x, n)
	case *details.Details:
		return splitDetails(x, n)
	case *token.Legacy:
		return splitLegacy(x, n)
	}
	return nil, false
}

func splitPlant(p *plant.Plant, n int) (*Item, bool) {
	if p.Count <= n {
		return nil, false
	}
	p.Count -= n
	return New(&plant.Plant{Type: p.Type, Count: n}), true
}

func splitWallet(w *money.Wallet, n int) (*Item, bool) {
	// Not `<=`: allow empty wallet.
	if w.Money < n {
		return nil, false
	}
	w.Money -= n
	return New(&money.Cash{Money: n}), true

}

func splitCash(c *money.Cash, n int) (*Item, bool) {
	if c.Money <= n {
		return nil, false
	}
	c.Money -= n
	return New(&money.Cash{Money: n}), true
}

func splitDetails(d *details.Details, n int) (*Item, bool) {
	if d.Count <= n {
		return nil, false
	}
	d.Count -= n
	return New(&details.Details{Count: n}), true
}

func splitLegacy(l *token.Legacy, n int) (*Item, bool) {
	if l.Count <= n {
		return nil, false
	}
	l.Count -= n
	return New(&token.Legacy{Count: n}), true
}
