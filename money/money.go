package money

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

const Currency = "₽"

var (
	ErrNoMoney  = errors.New("insufficient money")
	ErrBadMoney = errors.New("incorrect amount of money")
)

// Cash represents some amount of money.
type Cash struct {
	Money int
}

// NewCash returns a random amount of cash.
func NewCash() *Cash {
	n := int(math.Abs(5000 + 2500*rand.NormFloat64()))
	return &Cash{Money: n}
}

// Spend implements the Spender interface.
func (c *Cash) Spend(n int) bool {
	if c.Money < n {
		return false
	}
	c.Money -= n
	return true
}

func (c Cash) String() string {
	return fmt.Sprintf("💵 Наличные (%d %s)", c.Money, Currency)
}

// Wallet is used to safely store money.
type Wallet struct {
	Money int
}

// NewWallet returns a wallet with a random amount of money.
func NewWallet() *Wallet {
	n := int(math.Abs(10000 + 5000*rand.NormFloat64()))
	return &Wallet{Money: n}
}

// Spend implements the Spender interface.
func (w *Wallet) Spend(n int) bool {
	if w.Money < n {
		return false
	}
	w.Money -= n
	return true
}

func (w Wallet) String() string {
	return fmt.Sprintf("💰 Кошелёк (%d %s)", w.Money, Currency)
}

type Transfer struct {
	Money   int
	Comment string
}

func (t Transfer) String() string {
	return fmt.Sprintf("💳 Перевод (%d %s, %s)", t.Money, Currency, t.Comment)
}
