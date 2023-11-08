package actions

import "nechego/item"

const inventoryCapacity = 20

func fullInventory(i *item.Set) bool {
	return i.Count() >= inventoryCapacity
}
