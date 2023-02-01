package item

import (
	"testing"
)

func TestRandom(t *testing.T) {
	const sample = 100_000
	items := make([]*Item, 0, sample)
	for i := 0; i < sample; i++ {
		items = append(items, Random())
	}
	types := map[Type]int{}
	for _, x := range items {
		if !x.Transferable {
			if x.Type != TypePair && x.Type != TypeEblan {
				t.Errorf("untransferable item of type %v", x.Type)
			}
		}
		switch x.Type {
		case TypeAdmin, TypeEblan, TypePair:
			if x.Expire.IsZero() {
				t.Errorf("zero expire of daily token %v", x.Type)
			}
		}
		if x.Type == TypeUnknown {
			t.Errorf("%v: item type unknown", x)
		}
		typ := TypeOf(x.Value)
		if typ != x.Type {
			t.Errorf("t == %v, want %v", typ, x.Type)
		}
		types[x.Type]++
	}
	wantTypes := []Type{
		TypeFish,
		TypeFood,
		TypeCash,
		TypeWallet,
		TypeFishingRod,
		TypeDetails,
		TypeThread,
		TypePet,
		TypeKnife,
		TypePhone,
		TypeDice,
		TypeAdmin,
	}
	for _, typ := range wantTypes {
		if types[typ] == 0 {
			t.Errorf("type %v not covered", typ)
		}
	}
}
