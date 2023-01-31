package item

import (
	"nechego/details"
	"nechego/fishing"
	"nechego/food"
	"nechego/money"
	"nechego/pets"
	"nechego/phone"
	"nechego/token"
	"nechego/tools"
)

func New(x any) *Item {
	i := &Item{
		Type:         TypeOf(x),
		Transferable: true,
		Value:        x,
	}
	if i.Type == TypeEblan || i.Type == TypePair {
		i.Transferable = false
	}
	return i
}

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
