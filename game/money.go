package game

import (
	"errors"
	"fmt"
	"nechego/item"
	"nechego/money"
)

var ErrNoMoney = errors.New("no money")

// Balance represents the user's money.
type Balance struct {
	inventory *item.Set
}

func (u *User) Balance() *Balance {
	return &Balance{u.Inventory}
}

// Rich is true if there is a lot of money on the balance.
func (b *Balance) Rich() bool {
	return b.Total() > 1_000_000
}

// Poor is true if there is almost no money on the balance.
func (b *Balance) Poor() bool {
	return b.Total() < 3000
}

// Total returns the aggregated amount of money on the balance.
func (b *Balance) Total() int {
	t := 0
	for _, x := range b.inventory.List() {
		switch v := x.Value.(type) {
		case *money.Cash:
			t += v.Money
		case *money.Wallet:
			t += v.Money
		}
	}
	return t
}

// Spend subtracts the specified amount of money from the balance and
// returns true. If there is not enough money on the balance, returns
// false.
func (b *Balance) Spend(n int) bool {
	if n < 0 {
		panic(fmt.Sprintf("cannot spend %d money", n))
	}
	for _, x := range b.inventory.List() {
		if x.Type != item.TypeCash && x.Type != item.TypeWallet {
			continue
		}
		if v, ok := x.Value.(Spender); ok && v.Spend(n) {
			return true
		}
	}
	return false
}

// Add adds a cash item of the specified value to the inventory.
func (b *Balance) Add(n int) {
	if n < 0 {
		panic(fmt.Sprintf("cannot add %d money", n))
	}
	if n == 0 {
		return
	}
	b.inventory.Add(item.New(&money.Cash{Money: n}))
}

// Cashout adds a cash item of the specified value to the inventory if
// there is enough money to do so.
func (b *Balance) Cashout(n int) error {
	if n <= 0 {
		return money.ErrBadMoney
	}
	if !b.Spend(n) {
		return money.ErrNoMoney
	}
	b.Add(n)
	return nil
}
