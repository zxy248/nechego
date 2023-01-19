package game

import (
	"nechego/fishing"
	"nechego/food"
	"nechego/item"
	"nechego/pets"
	"nechego/tools"
)

type craftFunc func(recipe []*item.Item) (result []*item.Item, ok bool)

func (u *User) Craft(recipe []*item.Item) (result []*item.Item, ok bool) {
	funcs := [...]craftFunc{
		craftRod,
		craftMeat,
	}
	for _, f := range funcs {
		if result, ok = f(recipe); ok {
			for _, x := range recipe {
				u.Inventory.Remove(x)
			}
			for _, x := range result {
				u.Inventory.Add(x)
			}
			return result, true
		}
	}
	return nil, false
}

func craftRod(recipe []*item.Item) (result []*item.Item, ok bool) {
	if len(recipe) != 2 {
		return nil, false
	}
	r0, ok := recipe[0].Value.(*fishing.Rod)
	if !ok {
		return nil, false
	}
	r1, ok := recipe[1].Value.(*fishing.Rod)
	if !ok {
		return nil, false
	}
	rod := &item.Item{
		Type:         item.TypeFishingRod,
		Transferable: true,
		Value: &fishing.Rod{
			Quality:    (r0.Quality+r1.Quality)/2 + 0.1,
			Durability: (r0.Durability + r1.Durability) / 2,
		},
	}
	return []*item.Item{rod}, true
}

func craftMeat(recipe []*item.Item) (result []*item.Item, ok bool) {
	if len(recipe) != 2 {
		return nil, false
	}
	knife, ok := recipe[0].Value.(*tools.Knife)
	if !ok {
		return nil, false
	}
	pet, ok := recipe[1].Value.(*pets.Pet)
	if !ok {
		return nil, false
	}
	return []*item.Item{
		{
			Type:         item.TypeMeat,
			Transferable: true,
			Value:        &food.Meat{Species: pet.Species},
		},
		{
			Type:         item.TypeKnife,
			Transferable: true,
			Value:        &tools.Knife{Durability: knife.Durability - 0.1},
		},
	}, true
}
