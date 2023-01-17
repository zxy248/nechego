package game

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"nechego/fishing"
	"nechego/food"
	"nechego/pets"
	"time"
)

type ItemType int

const (
	ItemTypeUnknown ItemType = iota
	ItemTypeEblanToken
	ItemTypeAdminToken
	ItemTypePairToken
	ItemTypeCash
	ItemTypeWallet
	ItemTypeCreditCard
	ItemTypeDebt
	ItemTypeFishingRod
	ItemTypeFish
	ItemTypePet
	ItemTypeDice
	ItemTypeFood
)

type Item struct {
	Type         ItemType
	Transferable bool
	Expire       time.Time
	Value        any
}

func (i *Item) UnmarshalJSON(data []byte) error {
	type ItemWrapper *Item

	var raw json.RawMessage
	wrapped := ItemWrapper(i) // prevent infinite recursion
	wrapped.Value = &raw
	if err := json.Unmarshal(data, wrapped); err != nil {
		return err
	}

	switch i.Type {
	case ItemTypeEblanToken:
		i.Value = &EblanToken{}
	case ItemTypeAdminToken:
		i.Value = &AdminToken{}
	case ItemTypePairToken:
		i.Value = &PairToken{}
	case ItemTypeCash:
		i.Value = &Cash{}
	case ItemTypeWallet:
		i.Value = &Wallet{}
	case ItemTypeCreditCard:
		i.Value = &CreditCard{}
	case ItemTypeDebt:
		i.Value = &Debt{}
	case ItemTypeFishingRod:
		i.Value = &FishingRod{}
	case ItemTypeFish:
		i.Value = &fishing.Fish{}
	case ItemTypePet:
		i.Value = &pets.Pet{}
	case ItemTypeDice:
		i.Value = &Dice{}
	case ItemTypeFood:
		i.Value = &food.Food{}
	default:
		panic(fmt.Errorf("unexpected item type %v", i.Type))
	}
	return json.Unmarshal(raw, i.Value)
}

func (i *Item) bad() bool {
	if rod, ok := i.Value.(*FishingRod); ok && rod.Durability < 0 {
		return true
	}
	if !i.Expire.IsZero() && time.Now().After(i.Expire) {
		return true
	}
	return false
}

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

func (it *Items) Exist(x *Item) bool {
	for _, v := range it.I {
		if v == x {
			return true
		}
	}
	return false
}

func (it *Items) ByKey(k int) (x *Item, ok bool) {
	x, ok = it.keys[k]
	if !ok || !it.Exist(x) {
		return nil, false
	}
	return x, true
}

func (it *Items) updateHotkeys() {
	it.keys = map[int]*Item{}
	for i, v := range it.I {
		it.keys[i] = v
	}
}

func (it *Items) normalize() []*Item {
	it.Filter(func(i *Item) bool { return !i.bad() })
	return it.I
}

func (it *Items) List() []*Item {
	// updates hotkeys; only for public use
	it.normalize()
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

func (it *Items) Retain(n int) {
	if len(it.I) > n {
		it.I = it.I[len(it.I)-n:]
	}
}

func (it *Items) Random() (x *Item, ok bool) {
	items := it.normalize()
	if len(items) == 0 {
		return nil, false
	}
	return items[rand.Intn(len(items))], true
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

const (
	InventorySize = 10
	InventoryCap  = 17
)

func (it *Items) Count() int {
	return len(it.I)
}

type SetNamer interface {
	SetName(s string) bool
}

func (i *Item) SetName(s string) bool {
	if x, ok := i.Value.(SetNamer); ok {
		return x.SetName(s)
	}
	return false
}

func randomItem() *Item {
	items := []*Item{
		{Type: ItemTypeFishingRod, Value: NewFishingRod()},
		{Type: ItemTypeFish, Value: fishing.RandomFish()},
		{Type: ItemTypePet, Value: pets.Random()},
		{Type: ItemTypeFood, Value: food.Random()},
		{Type: ItemTypeWallet, Value: &Wallet{Money: int(math.Abs(rand.NormFloat64() * 10000))}},
		{Type: ItemTypeCash, Value: &Cash{Money: int(math.Abs(rand.NormFloat64() * 2000))}},
	}
	if rand.Float64() < 0.50 {
		items = append(items, &Item{Type: ItemTypeDice, Value: &Dice{}})
	}
	if rand.Float64() < 0.02 {
		items = append(items, &Item{Type: ItemTypeAdminToken, Expire: tomorrow(), Value: &AdminToken{}})
	}
	i := items[rand.Intn(len(items))]
	i.Transferable = true
	return i
}
