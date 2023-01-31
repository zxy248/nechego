package item

import (
	"math/rand"
)

type Items struct {
	I    []*Item
	keys map[int]*Item
}

func NewItems() *Items {
	return &Items{I: []*Item{}}
}

func (it *Items) Add(x *Item) {
	it.I = append(it.I, x)
}

func (it *Items) Remove(x *Item) bool {
	for i, v := range it.I {
		if v == x {
			it.I[i] = it.I[len(it.I)-1]
			it.I = it.I[:len(it.I)-1]
			return true
		}
	}
	return false
}

func (it *Items) Contain(x *Item) bool {
	for _, v := range it.I {
		if v == x {
			return true
		}
	}
	return false
}

func (it *Items) ByKey(k int) (x *Item, ok bool) {
	x, ok = it.keys[k]
	if !ok || !it.Contain(x) {
		return nil, false
	}
	return x, true
}

func (it *Items) ByType(t Type) (x *Item, ok bool) {
	for _, x := range it.List() {
		if x.Type == t {
			return x, true
		}
	}
	return nil, false
}

func (it *Items) updateHotkeys() {
	it.keys = map[int]*Item{}
	for i, v := range it.I {
		it.keys[i] = v
	}
}

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

func (it *Items) List() []*Item {
	it.Filter(integral)
	return it.I
}

func (it *Items) HkList() []*Item {
	// Updates hotkeys. Do not use this function internally.
	it.Filter(integral)
	it.updateHotkeys()
	return it.I
}

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

func (it *Items) Trim(n int) {
	if len(it.I) > n {
		it.I = it.I[len(it.I)-n:]
	}
}

func (it *Items) Random() (x *Item, ok bool) {
	items := it.List()
	if len(items) == 0 {
		return nil, false
	}
	return items[rand.Intn(len(items))], true
}

func (it *Items) Count() int {
	return len(it.I)
}

func (it *Items) PushFront(i []*Item) {
	it.I = append(i, it.I...)
}
