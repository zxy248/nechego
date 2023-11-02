package item

import (
	"fmt"
	"nechego/details"
	"nechego/farm/plant"
	"nechego/fishing"
	"nechego/food"
	"nechego/money"
	"nechego/pets"
	"nechego/token"
	"nechego/tools"
)

// The dynamic type of an Item corresponding to the actual type of its
// underlying value.
type Type int

const (
	TypeUnknown Type = iota
	TypeEblan
	TypeAdmin
	TypePair
	TypeCash
	TypeWallet
	TypeFishingRod
	TypeFish
	TypePet
	TypeFood
	TypeKnife
	TypeMeat
	TypeDetails
	TypeThread
	TypeFishingNet
	TypePlant
	TypeLegacy
)

// TypeOf returns a Type corresponding to the actual type of x.
func TypeOf(x any) Type {
	switch x.(type) {
	case *token.Eblan:
		return TypeEblan
	case *token.Admin:
		return TypeAdmin
	case *token.Pair:
		return TypePair
	case *money.Cash:
		return TypeCash
	case *money.Wallet:
		return TypeWallet
	case *fishing.Rod:
		return TypeFishingRod
	case *fishing.Fish:
		return TypeFish
	case *pets.Pet:
		return TypePet
	case *food.Food:
		return TypeFood
	case *tools.Knife:
		return TypeKnife
	case *food.Meat:
		return TypeMeat
	case *details.Details:
		return TypeDetails
	case *details.Thread:
		return TypeThread
	case *fishing.Net:
		return TypeFishingNet
	case *plant.Plant:
		return TypePlant
	case *token.Legacy:
		return TypeLegacy
	default:
		return TypeUnknown
	}
}

// ValueOf returns the dynamic value of the specified type.
// Panics if the type t is not supported.
func ValueOf(t Type) any {
	switch t {
	case TypeEblan:
		return &token.Eblan{}
	case TypeAdmin:
		return &token.Admin{}
	case TypePair:
		return &token.Pair{}
	case TypeCash:
		return &money.Cash{}
	case TypeWallet:
		return &money.Wallet{}
	case TypeFishingRod:
		return &fishing.Rod{}
	case TypeFish:
		return &fishing.Fish{}
	case TypePet:
		return &pets.Pet{}
	case TypeFood:
		return &food.Food{}
	case TypeKnife:
		return &tools.Knife{}
	case TypeMeat:
		return &food.Meat{}
	case TypeDetails:
		return &details.Details{}
	case TypeThread:
		return &details.Thread{}
	case TypeFishingNet:
		return &fishing.Net{}
	case TypePlant:
		return &plant.Plant{}
	case TypeLegacy:
		return &token.Legacy{}
	default:
		panic(fmt.Sprintf("unexpected item type %v", t))
	}
}
