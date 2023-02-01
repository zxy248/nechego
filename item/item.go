package item

import (
	"encoding/json"
	"math"
	"math/rand"
	"nechego/dates"
	"nechego/details"
	"nechego/fishing"
	"nechego/food"
	"nechego/money"
	"nechego/pets"
	"nechego/phone"
	"nechego/token"
	"nechego/tools"
	"time"
)

// Item represents an item in the world.
type Item struct {
	Type         Type      // Type of the underlying item value.
	Transferable bool      // Transferable is true if the item can be transfered.
	Expire       time.Time // Expire specifies the time after which the item is gone.
	Value        any       // Value of the actual object.
}

// New wraps an item value in Item.
func New(x any) *Item {
	i := &Item{
		Type:         TypeOf(x),
		Transferable: true,
		Value:        x,
	}
	switch i.Type {
	case TypeEblan, TypePair:
		i.Transferable = false
	}
	switch i.Type {
	case TypeEblan, TypePair, TypeAdmin:
		i.Expire = dates.Tomorrow()
	}
	return i
}

// UnmarshalJSON decodes the Item from its textual representation.
func (i *Item) UnmarshalJSON(data []byte) error {
	// Necessary to prevent infinite recursion.
	type ItemWrapper *Item

	// Value should be decoded after Type is known.
	var raw json.RawMessage
	wrapped := ItemWrapper(i)
	wrapped.Value = &raw
	if err := json.Unmarshal(data, wrapped); err != nil {
		return err
	}

	// Now the dynamic type is accessible; assign and unmarshal
	// the underliying object.
	wrapped.Value = ValueOf(i.Type)
	return json.Unmarshal(raw, i.Value)
}

type SetNamer interface {
	SetName(s string) bool
}

// SetName sets the name of the underlying object if it implements the
// SetNamer interface.
func (i *Item) SetName(s string) bool {
	if x, ok := i.Value.(SetNamer); ok {
		return x.SetName(s)
	}
	return false
}

// Random returns a random item.
func Random() *Item {
	common := []any{
		fishing.RandomFish(),
		food.Random(),
		&money.Cash{Money: int(math.Abs(rand.NormFloat64() * 3000))},
	}
	uncommon := []any{
		&money.Wallet{Money: int(math.Abs(rand.NormFloat64() * 10000))},
		fishing.NewRod(),
		details.Random(),
		&details.Thread{},
	}
	rare := []any{
		pets.Random(),
		&tools.Knife{Durability: 0.8 + 0.2*rand.Float64()},
	}
	epic := []any{
		phone.NewPhone(),
		&token.Dice{},
	}
	legendary := []any{
		&token.Admin{},
	}
	table := []struct {
		threshold float64
		list      []any
	}{
		{1.0, common},
		{0.5, uncommon},
		{0.25, rare},
		{0.12, epic},
		{0.02, legendary},
	}
	items := []any{}
	for _, x := range table {
		if rand.Float64() < x.threshold {
			items = append(items, x.list...)
		}
	}
	return New(items[rand.Intn(len(items))])
}

// integral returns true if the item is OK, and returns false if the
// item should be removed.
func integral(i *Item) bool {
	switch x := i.Value.(type) {
	case *fishing.Rod:
		if x.Durability < 0 {
			return false
		}
	case *money.Cash:
		if x.Money == 0 {
			return false
		}
	case *tools.Knife:
		if x.Durability < 0 {
			return false
		}
	case *details.Details:
		if x.Count == 0 {
			return false
		}
	case *fishing.Net:
		if x.Count() == 0 && x.Broken() {
			return false
		}
	}
	if !i.Expire.IsZero() && time.Now().After(i.Expire) {
		return false
	}
	return true
}
