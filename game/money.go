package game

import (
	"fmt"
	"nechego/fishing"
	"nechego/item"
	"nechego/money"
)

const (
	RichThreshold = 1_000_000
	PoorThreshold = 3_000
)

func (u *User) Rich() bool {
	return u.Total() >= RichThreshold
}

func (u *User) Poor() bool {
	return u.Total() < PoorThreshold
}

func (u *User) SpendMoney(n int) bool {
	if n < 0 {
		panic(fmt.Errorf("cannot spend %v", n))
	}
	u.Stack()
	return u.spendWallet(n) || u.spendCash(n)
}

func (u *User) spendWallet(n int) bool {
	w, ok := u.Wallet()
	if !ok {
		return false
	}
	if w.Money < n {
		return false
	}
	w.Money -= n
	return true
}

func (u *User) spendCash(n int) bool {
	c, ok := u.Cash()
	if !ok {
		return false
	}
	if c.Money < n {
		return false
	}
	c.Money -= n
	return true
}

func (u *User) AddMoney(n int) {
	u.Inventory.Add(&item.Item{
		Type:         item.TypeCash,
		Transferable: true,
		Value:        &money.Cash{Money: n},
	})
}

func (u *User) Stack() bool {
	t := 0
	u.Inventory.Filter(func(i *item.Item) bool {
		switch x := i.Value.(type) {
		case *money.Cash:
			t += x.Money
			return false
		case *money.Wallet:
			t += x.Money
			x.Money = 0
		}
		return true
	})
	if t == 0 {
		return false
	}
	wallet, ok := u.Wallet()
	if !ok {
		u.AddMoney(t)
		return true
	}
	wallet.Money += t
	return true
}

func (u *User) Cashout(n int) error {
	if n <= 0 {
		return money.ErrBadMoney
	}
	if !u.SpendMoney(n) {
		return money.ErrNoMoney
	}
	u.AddMoney(n)
	return nil
}

func (u *User) Total() int {
	t := 0
	for _, v := range u.Inventory.Normal() {
		switch x := v.Value.(type) {
		case *money.Cash:
			t += x.Money
		case *money.Wallet:
			t += x.Money
		case *money.CreditCard:
			t += x.Money
		case *money.Debt:
			t -= x.Money
		}
	}
	return t
}

func (u *User) InDebt() bool {
	for _, v := range u.Inventory.Normal() {
		if _, ok := v.Value.(*money.Debt); ok {
			return true
		}
	}
	return false
}

func (u *User) Sell(i *item.Item) (profit int, ok bool) {
	if !i.Transferable {
		return 0, false
	}
	if ok = u.Inventory.Remove(i); !ok {
		return 0, false
	}
	switch x := i.Value.(type) {
	case *fishing.Fish:
		n := int(x.Price())
		u.AddMoney(n)
		return n, true
	default:
		// cannot sell; return item back
		u.Inventory.Add(i)
	}
	return 0, false
}
