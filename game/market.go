package game

import (
	"math/rand"
	"nechego/fishing"
	"nechego/pets"
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
	fishPrice := int(fish.Price() * (0.5 + rand.Float64()))
	products := []*Product{
		{2500 + rand.Intn(5000), &Item{
			Type:         ItemTypeFishingRod,
			Transferable: true,
			Value:        NewFishingRod()}},
		{fishPrice, &Item{
			Type:         ItemTypeFish,
			Transferable: true,
			Value:        fish}},
		{100 + rand.Intn(50000), &Item{
			Type:         ItemTypePet,
			Transferable: true,
			Value:        pets.RandomPet()}},
		{500 + rand.Intn(4500), &Item{
			Type:         ItemTypeDice,
			Transferable: true,
			Value:        &Dice{}}},
	}
	if rand.Float64() < 0.25 {
		products = append(products, &Product{
			500000 + rand.Intn(1000000), &Item{
				Type:         ItemTypeAdminToken,
				Transferable: true,
				Value:        &AdminToken{}}})
	}
	m.Add(products[rand.Intn(len(products))])
	const maxitems = 10
	if len(m.P) > maxitems {
		m.P = m.P[len(m.P)-maxitems:]
	}
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
