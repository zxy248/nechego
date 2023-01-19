package game

import (
	"errors"
	"fmt"
	"math/rand"
	"nechego/fishing"
	"nechego/item"
	"nechego/money"
	"nechego/valid"
	"strings"
)

var ErrNoKey = errors.New("key not found")

type Product struct {
	Price int
	Item  *item.Item
}

type Market struct {
	P    []*Product
	Name string
	keys map[int]*Product
}

func NewMarket() *Market {
	return &Market{P: []*Product{}}
}

func (m *Market) Refill() {
	product := randomProduct()
	m.Add(product)
	const maxitems = 10
	if len(m.P) > maxitems {
		m.P = m.P[len(m.P)-maxitems:]
	}
}

func randomProduct() *Product {
	p := 0
	i := item.Random()
	switch i.Type {
	case item.TypeFishingRod:
		p = 2500 + rand.Intn(7500)
	case item.TypeFish:
		f := i.Value.(*fishing.Fish)
		p = int(f.Price() * (0.5 + 1.5*rand.Float64()))
	case item.TypePet:
		p = 500 + rand.Intn(99500)
	case item.TypeDice:
		p = 5000 + rand.Intn(25000)
	case item.TypeFood:
		p = 250 + rand.Intn(1750)
	case item.TypeAdmin:
		p = 500_000 + rand.Intn(4_500_000)
	case item.TypeKnife:
		p = 1000 + rand.Intn(9000)
	default:
		return randomProduct()
	}
	return &Product{p, i}
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

func (m *Market) SetName(s string) bool {
	if !valid.Name(s) {
		return false
	}
	m.Name = strings.Title(s)
	return true
}

func (m *Market) String() string {
	s := "🏪 Магазин"
	if m.Name != "" {
		s += fmt.Sprintf(` «%s»`, m.Name)
	}
	return s
}

func (u *User) Buy(m *Market, key int) (*Product, error) {
	p, ok := m.keys[key]
	if !ok {
		return nil, ErrNoKey
	}
	if !u.SpendMoney(p.Price) {
		return nil, money.ErrNoMoney
	}
	delete(m.keys, key)
	for i, v := range m.P {
		if v == p {
			m.P[i] = m.P[len(m.P)-1]
			m.P = m.P[:len(m.P)-1]
		}
	}
	u.Inventory.Add(p.Item)
	return p, nil
}
