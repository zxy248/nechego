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
	var p float64
	i := item.Random()
	switch x := i.Value.(type) {
	case *fishing.Rod:
		p = 5000
	case *fishing.Fish:
		p = x.Price()
	case *pets.Pet:
		switch q := x.Species.Quality(); q {
		case pets.Common:
			p = 5e3
		case pets.Rare:
			p = 5e4
		case pets.Exotic:
			p = 5e5
		case pets.Legendary:
			p = 5e6
		}
	case *food.Food:
		p = x.Price()
	case *tools.Knife:
		p = 2000
	case *details.Details:
		p = float64(x.Count) * 150
	case *plant.Plant:
		p = float64(x.Count) * x.Price()
	default:
		return randomProduct()
	}
	q := int(math.Abs(p + 0.16*p*rand.NormFloat64()))
	return &Product{q, i}
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
