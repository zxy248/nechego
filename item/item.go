package item

import (
	"encoding/json"
	"math/rand"
	"nechego/dates"
	"nechego/details"
	"nechego/farm/plant"
	"nechego/fishing"
	"nechego/food"
	"nechego/money"
	"nechego/pets"
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

// New returns an Item of the value x.
func New(x any) *Item {
	i := &Item{
		Type:         TypeOf(x),
		Transferable: true,
		Value:        x,
	}
	switch i.Type {
	case TypeEblan, TypePair, TypeAdmin:
		i.Transferable = false
		i.Expire = dates.Tomorrow()
	}
	return i
}

// UnmarshalJSON implements the json.Unmarshaler interface.
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

// Random returns a random item.
func Random() *Item {
	pool := map[float64][]any{
		1.0: {
			fishing.RandomFish(),
			food.Random(),
			money.NewCash(),
			plant.Random(),
		},
		0.5: {
			money.NewWallet(),
			fishing.NewRod(),
		},
		0.37: {
			details.Random(),
			pets.Random(),
		},
		0.25: {
			&details.Thread{},
		},
		0.12: {
			tools.NewKnife(),
		},
	}
	i := []any{}
	for threshold, list := range pool {
		if rand.Float64() < threshold {
			i = append(i, list...)
		}
	}
	return New(i[rand.Intn(len(i))])
}

// integral returns true if the item should be kept, and returns false
// if the item should be removed.
func integral(i *Item) bool {
	switch x := i.Value.(type) {
	case *fishing.Rod:
		if x.Broken() {
			return false
		}
	case *money.Cash:
		if x.Money == 0 {
			return false
		}
	case *tools.Knife:
		if x.Broken() {
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
	case *plant.Plant:
		if x.Count == 0 {
			return false
		}
	}
	if !i.Expire.IsZero() && time.Now().After(i.Expire) {
		return false
	}
	return true
}
