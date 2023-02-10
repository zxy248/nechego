package game

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"nechego/details"
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

// Product is an item with price to be sold on the market.
type Product struct {
	Price int
	Item  *item.Item
}

// Market represents a place where a user can buy products.
type Market struct {
	P     []*Product       // P is a list of products on sale.
	Name  string           // Name of the market.
	Shift Shift            // Work shift.
	keys  map[int]*Product // keys for product selection.
}

// NewMarket returns a new Market with no products on sale.
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
			panic(fmt.Sprintf("unexpected pet type %d", q))
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
	case *details.Details:
		p = price(5000, 2500)
	case *details.Thread:
		p = price(5000, 2500)
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

// SetName sets the market's name.
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

// Buy removes the product specified by key from the market and adds
// it to the user's inventory if there is enough money on the balance.
func (u *User) Buy(w *World, key int) (*Product, error) {
	market := w.Market

	product, ok := market.keys[key]
	if !ok {
		return nil, ErrNoKey
	}
	if !u.Balance().Spend(product.Price) {
		return nil, money.ErrNoMoney
	}

	// The purchase is commited.
	delete(market.keys, key)
	for i, v := range market.P {
		if v == product {
			market.P[i] = market.P[len(market.P)-1]
			market.P = market.P[:len(market.P)-1]
		}
	}
	u.Inventory.Add(product.Item)

	if earn := product.Price / 3; earn > 0 {
		// The market worker earns a wage.
		if id, ok := market.Shift.Worker(); ok {
			worker := w.UserByID(id)
			worker.Funds.Add("рохля", item.New(&money.Cash{Money: earn}))
		}
		// The strongest player takes a tax.
		if top := w.SortedUsers(ByStrength); len(top) > 0 {
			strongest := top[0]
			strongest.Funds.Add("налог", item.New(&money.Cash{Money: earn}))
		}
	}
	return product, nil
}

// Sell removes the specified item from the inventory and adds money
// if the item can be sold.
func (u *User) Sell(w *World, i *item.Item) (profit int, ok bool) {
	if !i.Transferable {
		return 0, false
	}

	// The item will be either sold or returned back to the inventory.
	if !u.Inventory.Remove(i) {
		return 0, false
	}

	switch x := i.Value.(type) {
	case *fishing.Fish:
		profit = int(x.Price())
	case *details.Details:
		profit = 100 * x.Count
	default:
		// Item of this type cannot be sold; return it back.
		u.Inventory.Add(i)
		return 0, false
	}
	// The sale is commited.
	u.Balance().Add(profit)

	// The strongest player takes a tax.
	if top, earn := w.SortedUsers(ByStrength), profit/10; len(top) > 0 && earn > 0 {
		strongest := top[0]
		strongest.Funds.Add("налог", item.New(&money.Cash{Money: earn}))
	}
	return profit, true
}
