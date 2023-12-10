package handlers

import "nechego/item"

const InventoryCapacity = 20

func FullInventory(i *item.Set) bool {
	return i.Count() >= InventoryCapacity
}

func GetItems(s *item.Set, ks []int) []*item.Item {
	var items []*item.Item
	seen := map[*item.Item]bool{}
	for _, k := range ks {
		x, ok := s.ByKey(k)
		if !ok || seen[x] {
			break
		}
		seen[x] = true
		items = append(items, x)
	}
	return items
}

func MoveItems(dst, src *item.Set, items []*item.Item) (moved []*item.Item, bad *item.Item) {
	for _, x := range items {
		if !src.Move(dst, x) {
			return moved, x
		}
		moved = append(moved, x)
	}
	return
}
