package game

import (
	"fmt"
	"nechego/item"
	"nechego/modifier"
	"nechego/money"
)

// Balance represents the user's money.
type Balance struct {
	inventory *item.Items
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
		panic(fmt.Errorf("cannot spend %v money", n))
	}
	b.Stack()
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
		panic(fmt.Errorf("cannot add %v money", n))
	}
	b.inventory.Add(&item.Item{
		Type:         item.TypeCash,
		Transferable: true,
		Value:        &money.Cash{Money: n},
	})
}

// Stack aggregates all money found in the inventory in a single slot.
func (b *Balance) Stack() {
	total := 0
	var wallet *money.Wallet
	b.inventory.Filter(func(i *item.Item) bool {
		switch x := i.Value.(type) {
		case *money.Cash:
			total += x.Money
			// Don't keep zero value cash.
			return false
		case *money.Wallet:
			// Stack money to the first wallet found.
			if wallet == nil {
				wallet = x
			}
			total += x.Money
			x.Money = 0
		}
		return true
	})
	if total == 0 {
		return
	}
	if wallet == nil {
		b.Add(total)
		return
	}
	wallet.Money += total
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

// Mod returns a modifier corresponding to the amount of money on the balance.
func (b *Balance) Mod() (m *modifier.Mod, ok bool) {
	if b.Rich() {
		return &modifier.Mod{
			Emoji:       "🎩",
			Multiplier:  +0.05,
			Description: "Вы богаты.",
		}, true
	}
	if b.Poor() {
		return &modifier.Mod{
			Emoji:       "🗑️",
			Multiplier:  -0.05,
			Description: "Вы бедны.",
		}, true
	}
	return nil, false
}
