package recipes

import (
	"nechego/details"
	"nechego/fishing"
	"nechego/food"
	"nechego/item"
	"nechego/pets"
	"nechego/tools"
)

type Recipe func([]*item.Item) (result []*item.Item, ok bool)

func UpgradeFishingRod(recipe []*item.Item) (result []*item.Item, ok bool) {
	if !(template{item.TypeFishingRod, item.TypeFishingRod}.match(recipe)) {
		return nil, false
	}
	r := recipe[0].Value.(*fishing.Rod)
	q := recipe[1].Value.(*fishing.Rod)
	if r.Level != q.Level {
		return nil, false
	}
	r.Level++
	r.Durability = (r.Durability + q.Durability) / 2
	return []*item.Item{{Type: item.TypeFishingRod, Value: r}}, true
}

func RepairFishingRod(recipe []*item.Item) (result []*item.Item, ok bool) {
	if !(template{item.TypeDetails, item.TypeFishingRod}.match(recipe)) {
		return nil, false
	}
	d := recipe[0].Value.(*details.Details)
	r := recipe[1].Value.(*fishing.Rod)
	need := int((1 - r.Durability) * 100)
	resource := d.Count
	if need < resource {
		resource = need
	}
	d.Spend(resource)
	r.Durability += float64(resource) * 0.01
	return []*item.Item{
		{Type: item.TypeFishingRod, Value: r},
		{Type: item.TypeDetails, Value: d},
	}, true
}

func DustFishingRod(recipe []*item.Item) (result []*item.Item, ok bool) {
	if !(template{item.TypeFishingRod}.match(recipe)) {
		return nil, false
	}
	r := recipe[0].Value.(*fishing.Rod)
	d := &details.Details{Count: r.Level * 10}
	return []*item.Item{{Type: item.TypeDetails, Value: d}}, true
}

func DustPhone(recipe []*item.Item) (result []*item.Item, ok bool) {
	if !(template{item.TypePhone}.match(recipe)) {
		return nil, false
	}
	d := &details.Details{Count: 50}
	return []*item.Item{{Type: item.TypePhone, Value: d}}, true
}

func MakeMeat(recipe []*item.Item) (result []*item.Item, ok bool) {
	if !(template{item.TypeKnife, item.TypePet}.match(recipe)) {
		return nil, false
	}
	knife := recipe[0].Value.(*tools.Knife)
	knife.Durability -= 0.1
	pet := recipe[1].Value.(*pets.Pet)
	meat := &food.Meat{Species: pet.Species}
	return []*item.Item{
		{Type: item.TypeMeat, Value: meat},
		{Type: item.TypeKnife, Value: knife},
	}, true
}
