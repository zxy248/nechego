package game

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"nechego/details"
	"nechego/farm/plant"
	"nechego/fishing"
	"nechego/food"
	"nechego/item"
	"nechego/money"
	"nechego/pets"
	"nechego/tools"
)

var ErrNoKey = errors.New("key not found")

// Product is an item with price to be sold on the market.
type Product struct {
	Price int
	Item  *item.Item
}

// Market represents a place where a user can buy products.
type Market struct {
	P         []*Product       // P is a list of products on sale.
	Name      string           // Name of the market.
	Shift     Shift            // Work shift.
	PriceList *PriceList       // Dynamic plant prices.
	keys      map[int]*Product // keys for product selection.

	OnBuy  func(*User, *Product)        `json:"-"`
	OnSell func(*User, *item.Item, int) `json:"-"`
}

// NewMarket returns a new Market with no products on sale.
func NewMarket() *Market {
	return &Market{
		P:         []*Product{},
		PriceList: NewPriceList(),
	}
}

// Refill adds a new random product to the market.
// If the number of products at the market would exceed a threshold,
// older products will be removed.
func (m *Market) Refill() {
	m.Add(randomProduct())
}

// Trim retains the last n market products, removing the preceding
// ones if necessary.
func (m *Market) Trim(n int) {
	if len(m.P) > n {
		m.P = m.P[len(m.P)-n:]
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
		p = normalize(x.Price() * (1.0 + 0.33*rand.NormFloat64()))
	case *pets.Pet:
		switch q := x.Species.Quality(); q {
		case pets.Common:
			p = price(3000, 1500)
		case pets.Rare:
			p = price(10000, 5000)
		case pets.Exotic:
			p = price(50000, 25000)
		case pets.Legendary:
			p = price(200000, 100000)
		default:
			panic(fmt.Sprintf("unexpected pet type %d", q))
		}
	case *food.Food:
		q := x.Price()
		p = price(q, q/3)
	case *tools.Knife:
		p = price(5000, 2000)
	case *details.Details:
		p = price(5000, 2500)
	case *plant.Plant:
		p = x.Count * price(1000, 500)
	default:
		// This type of item cannot be sold at the market.
		// Reroll.
		return randomProduct()
	}
	return &Product{p, i}
}

// Add adds a new product to the market.
func (m *Market) Add(p *Product) {
	m.P = append(m.P, p)
}

// Products returns a list of products at the market.
func (m *Market) Products() []*Product {
	m.keys = map[int]*Product{}
	for i, p := range m.P {
		m.keys[i] = p
	}
	return m.P
}

// String returns the textual representation of the Market.
func (m *Market) String() string {
	s := "🏪 Магазин"
	if m.Name != "" {
		s += fmt.Sprintf(" «%s»", m.Name)
	}
	return s
}

// Buy removes the product specified by key from the market and adds
// it to the user's inventory if there is enough money on the balance.
func (u *User) Buy(m *Market, key int) (*Product, error) {
	product, ok := m.keys[key]
	if !ok {
		return nil, ErrNoKey
	}
	if !u.Balance().Spend(product.Price) {
		return nil, money.ErrNoMoney
	}

	// The purchase is commited.
	delete(m.keys, key)
	for i, v := range m.P {
		if v == product {
			m.P[i] = m.P[len(m.P)-1]
			m.P = m.P[:len(m.P)-1]
		}
	}
	u.Inventory.Add(product.Item)
	if m.OnBuy != nil {
		m.OnBuy(u, product)
	}
	return product, nil
}

// Sell removes the specified item from the inventory and adds money
// if the item can be sold.
func (u *User) Sell(m *Market, i *item.Item) (profit int, ok bool) {
	profit, ok = m.PriceList.Price(i)
	if !ok {
		return 0, false
	}
	u.Inventory.Remove(i)
	u.Balance().Add(profit)
	if m.OnSell != nil {
		m.OnSell(u, i, profit)
	}
	return profit, true
}
