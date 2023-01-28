package item

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"nechego/dates"
	"nechego/fishing"
	"nechego/food"
	"nechego/money"
	"nechego/pets"
	"nechego/phone"
	"nechego/token"
	"nechego/tools"
	"time"
)

type Type int

const (
	TypeUnknown Type = iota
	TypeEblan
	TypeAdmin
	TypePair
	TypeCash
	TypeWallet
	TypeCreditCard // TODO: remove on wipe
	TypeDebt       // TODO: remove on wipe
	TypeFishingRod
	TypeFish
	TypePet
	TypeDice
	TypeFood
	TypeKnife
	TypeMeat
	TypePhone
)

type Item struct {
	Type         Type
	Transferable bool
	Expire       time.Time
	Value        any
}

func (i *Item) UnmarshalJSON(data []byte) error {
	// Necessary to prevent infinite recursion.
	type ItemWrapper *Item

	var raw json.RawMessage
	wrapped := ItemWrapper(i)
	wrapped.Value = &raw
	if err := json.Unmarshal(data, wrapped); err != nil {
		return err
	}

	switch i.Type {
	case TypeEblan:
		i.Value = &token.Eblan{}
	case TypeAdmin:
		i.Value = &token.Admin{}
	case TypePair:
		i.Value = &token.Pair{}
	case TypeCash:
		i.Value = &money.Cash{}
	case TypeWallet:
		i.Value = &money.Wallet{}
	case TypeCreditCard:
		// TODO: remove on wipe
	case TypeDebt:
		// TODO: remove on wipe
	case TypeFishingRod:
		i.Value = &fishing.Rod{}
	case TypeFish:
		i.Value = &fishing.Fish{}
	case TypePet:
		i.Value = &pets.Pet{}
	case TypeDice:
		i.Value = &token.Dice{}
	case TypeFood:
		i.Value = &food.Food{}
	case TypeKnife:
		i.Value = &tools.Knife{}
	case TypeMeat:
		i.Value = &food.Meat{}
	case TypePhone:
		i.Value = &phone.Phone{}
	default:
		panic(fmt.Errorf("unexpected item type %v", i.Type))
	}
	return json.Unmarshal(raw, i.Value)
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

func Random() *Item {
	common := []*Item{
		{Type: TypeFish, Value: fishing.RandomFish()},
		{Type: TypeFood, Value: food.Random()},
		{Type: TypeCash, Value: &money.Cash{Money: int(math.Abs(rand.NormFloat64() * 3000))}},
	}
	uncommon := []*Item{
		{Type: TypeWallet, Value: &money.Wallet{Money: int(math.Abs(rand.NormFloat64() * 10000))}},
		{Type: TypeFishingRod, Value: fishing.NewRod()},
	}
	rare := []*Item{
		{Type: TypePet, Value: pets.Random()},
		{Type: TypeKnife, Value: &tools.Knife{Durability: 0.8 + 0.2*rand.Float64()}},
	}
	epic := []*Item{
		{Type: TypePhone, Value: phone.NewPhone()},
		{Type: TypeDice, Value: &token.Dice{}},
	}
	legendary := []*Item{
		{Type: TypeAdmin, Expire: dates.Tomorrow(), Value: &token.Admin{}},
	}
	table := []struct {
		threshold float64
		list      []*Item
	}{
		{1.0, common},
		{0.5, uncommon},
		{0.25, rare},
		{0.12, epic},
		{0.02, legendary},
	}
	items := []*Item{}
	for _, x := range table {
		if rand.Float64() < x.threshold {
			items = append(items, x.list...)
		}
	}
	i := items[rand.Intn(len(items))]
	i.Transferable = true
	return i
}

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
	}
	if !i.Expire.IsZero() && time.Now().After(i.Expire) {
		return false
	}
	return true
}
