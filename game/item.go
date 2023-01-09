package game

import (
	"encoding/json"
	"fmt"
	"nechego/fishing"
	"time"
)

type ItemType int

const (
	ItemTypeUnknown ItemType = iota
	ItemTypeEblanToken
	ItemTypeAdminToken
	ItemTypePairToken
	ItemTypeWallet
	ItemTypeCreditCard
	ItemTypeDebt
	ItemTypeFishingRod
	ItemTypeFish
)

type Item struct {
	ID           int
	Type         ItemType
	Transferable bool
	Expire       time.Time
	Value        interface{}
}

func (i *Item) UnmarshalJSON(data []byte) error {
	type ItemWrapper *Item

	var raw json.RawMessage
	wrapped := ItemWrapper(i) // prevent infinite recursion
	wrapped.Value = &raw
	if err := json.Unmarshal(data, wrapped); err != nil {
		return err
	}

	switch i.Type {
	case ItemTypeEblanToken:
		i.Value = &EblanToken{}
	case ItemTypeAdminToken:
		i.Value = &AdminToken{}
	case ItemTypePairToken:
		i.Value = &PairToken{}
	case ItemTypeWallet:
		i.Value = &Wallet{}
	case ItemTypeCreditCard:
		i.Value = &CreditCard{}
	case ItemTypeDebt:
		i.Value = &Debt{}
	case ItemTypeFishingRod:
		i.Value = &FishingRod{}
	case ItemTypeFish:
		i.Value = &fishing.Fish{}
	default:
		panic(fmt.Errorf("unexpected item type %v", i.Type))
	}
	return json.Unmarshal(raw, i.Value)
}

func hotkeys(items []*Item) (map[int]*Item, []*Item) {
	m := map[int]*Item{}
	r := []*Item{}
	for i, v := range items {
		m[i] = v
		r = append(r, v)
	}
	return m, r
}
