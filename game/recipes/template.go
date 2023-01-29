package recipes

import "nechego/item"

// template verifies the number of ingredients and their type.
type template []item.Type

// match returns true if ingredients are isomorphic to the template.
func (t template) match(ingredients []*item.Item) bool {
	if len(t) != len(ingredients) {
		return false
	}
	for i, typ := range t {
		if typ != ingredients[i].Type {
			return false
		}
	}
	return true
}
