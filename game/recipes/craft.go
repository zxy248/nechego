package recipes

import "nechego/item"

// Craft matches ingredients with one of the recipes. If the
// appropriate recipe is found, removes ingredients from inventory and
// adds the result.
func Craft(inventory *item.Set, ingredients []*item.Item) (result []*item.Item, ok bool) {
	if hasDuplicates(ingredients) {
		return nil, false
	}
	list := [...]Recipe{
		UpgradeFishingRod,
		RepairFishingRod,
		DustFishingRod,
		DustPhone,
		MakeMeat,
		MakeFishingNet,
	}
	for _, craft := range list {
		if result, ok := craft(ingredients); ok {
			for _, x := range ingredients {
				inventory.Remove(x)
			}
			for _, x := range result {
				x.Transferable = true
				inventory.Add(x)
			}
			return result, true
		}
	}
	return nil, false
}

// hasDuplicates is true if some of the items in the list appear more
// than one time.
func hasDuplicates(items []*item.Item) bool {
	set := map[*item.Item]bool{}
	for _, x := range items {
		set[x] = true
	}
	if len(set) != len(items) {
		return true
	}
	return false
}
