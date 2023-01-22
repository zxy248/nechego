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
	TypeCreditCard
	TypeDebt
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
	type ItemWrapper *Item

	var raw json.RawMessage
	wrapped := ItemWrapper(i) // prevent infinite recursion
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
		i.Value = &money.CreditCard{}
	case TypeDebt:
		i.Value = &money.Debt{}
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
		{Type: TypeWallet, Value: &money.Wallet{Money: int(math.Abs(rand.NormFloat64() * 10000))}},
		{Type: TypeCash, Value: &money.Cash{Money: int(math.Abs(rand.NormFloat64() * 3000))}},
	}
	uncommon := []*Item{
		{Type: TypeFishingRod, Value: fishing.NewRod()},
		{Type: TypePet, Value: pets.Random()},
		{Type: TypePhone, Value: phone.NewPhone()},
	}
	rare := []*Item{
		{Type: TypeKnife, Value: &tools.Knife{Durability: 0.8 + 0.2*rand.Float64()}},
		{Type: TypeDice, Value: &token.Dice{}},
	}
	epic := []*Item{
		{Type: TypeAdmin, Expire: dates.Tomorrow(), Value: &token.Admin{}},
	}

	items := common
	if rand.Float64() < 0.5 {
		items = append(items, uncommon...)
	}
	if rand.Float64() < 0.25 {
		items = append(items, rare...)
	}
	if rand.Float64() < 0.02 {
		items = append(items, epic...)
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
