package item

import (
	"fmt"
	"nechego/details"
	"nechego/fishing"
	"nechego/food"
	"nechego/money"
	"nechego/pets"
	"nechego/phone"
	"nechego/token"
	"nechego/tools"
)

// Dynamic type of an item corresponding to the actual types.
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
	TypeDice
	TypeFood
	TypeKnife
	TypeMeat
	TypePhone
	TypeDetails
	TypeThread
	TypeFishingNet
)

// TypeOf returns Type of x.
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
	case *token.Dice:
		return TypeDice
	case *food.Food:
		return TypeFood
	case *tools.Knife:
		return TypeKnife
	case *food.Meat:
		return TypeMeat
	case *phone.Phone:
		return TypePhone
	case *details.Details:
		return TypeDetails
	case *details.Thread:
		return TypeThread
	case *fishing.Net:
		return TypeFishingNet
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
	case TypeDice:
		return &token.Dice{}
	case TypeFood:
		return &food.Food{}
	case TypeKnife:
		return &tools.Knife{}
	case TypeMeat:
		return &food.Meat{}
	case TypePhone:
		return &phone.Phone{}
	case TypeDetails:
		return &details.Details{}
	case TypeThread:
		return &details.Thread{}
	case TypeFishingNet:
		return &fishing.Net{}
	default:
		panic(fmt.Sprintf("unexpected item type %v", t))
	}
}
