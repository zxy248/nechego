package game

import (
	"math/rand"
	"nechego/fishing"
)

type Product struct {
	Price int
	Item  *Item
}

type Market struct {
	P    []*Product
	keys map[int]*Product
}

func NewMarket() *Market {
	return &Market{P: []*Product{}}
}

func (m *Market) Refill() {
	fish := fishing.RandomFish()
	fishPrice := int(float64(fish.Price()) * (0.5 + rand.Float64()))
	products := []*Product{
		{2000 + rand.Intn(2000), &Item{
			Type:         ItemTypeFishingRod,
			Transferable: true,
			Value:        NewFishingRod()}},
		{fishPrice, &Item{
			Type:         ItemTypeFish,
			Transferable: true,
			Value:        &fish}},
	}
	m.Add(products[rand.Intn(len(products))])
}

func (m *Market) Add(p *Product) {
	m.P = append(m.P, p)
}

func (m *Market) Products() []*Product {
	m.keys = map[int]*Product{}
	for i, p := range m.P {
		m.keys[i] = p
	}
	return m.P
}

func (u *User) Buy(m *Market, key int) (p *Product, ok bool) {
	p, ok = m.keys[key]
	if !ok {
		return nil, false
	}
	if ok := u.SpendMoney(p.Price); !ok {
		return nil, false
	}
	delete(m.keys, key)
	for i, v := range m.P {
		if v == p {
			m.P[i] = m.P[len(m.P)-1]
			m.P = m.P[:len(m.P)-1]
		}
	}
	u.Inventory.Add(p.Item)
	return p, true
}
