package game

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"nechego/fishing"
	"nechego/food"
	"nechego/item"
	"nechego/money"
	"nechego/pets"
	"nechego/phone"
	"nechego/token"
	"nechego/tools"
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

// Refill adds a new random product to the market.
// If the number of products at the market would exceed a threshold,
// older products will be removed.
func (m *Market) Refill() {
	const trim = 10
	product := randomProduct()
	m.Add(product)
	if len(m.P) > trim {
		m.P = m.P[len(m.P)-trim:]
	}
}

// randomProduct returns a random product that can be sold at the market.
func randomProduct() *Product {
	// normalize returns the absolute value of x as int.
	normalize := func(x float64) int { return int(math.Abs(x)) }

	// price returns a normally distributed positive int.
	price := func(mean, stddev float64) int {
		return normalize(mean + stddev*rand.NormFloat64())
	}

	var p int
	i := item.Random()
	switch x := i.Value.(type) {
	case *fishing.Rod:
		p = price(5000, 2500)
	case *fishing.Fish:
		p = normalize(x.Price() * (1.0 + 0.25*rand.NormFloat64()))
	case *pets.Pet:
		switch q := x.Species.Quality(); q {
		case pets.Common:
			p = price(3000, 1500)
		case pets.Rare:
			p = price(10000, 5000)
		case pets.Exotic:
			p = price(50000, 25000)
		case pets.Legendary:
			p = price(200_000, 100_000)
		default:
			panic(fmt.Errorf("unexpected pet type %d", q))
		}
	case *token.Dice:
		p = price(25000, 10000)
	case *food.Food:
		p = price(1000, 500)
	case *token.Admin:
		p = price(2_500_000, 1_000_000)
	case *tools.Knife:
		p = price(5000, 2000)
	case *phone.Phone:
		p = price(20000, 10000)
	default:
		// This type of item cannot be sold at the market.
		// Reroll.
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
