package money

import (
	"errors"
	"fmt"
)

const Currency = "₴"

var (
	ErrNoMoney  = errors.New("insufficient money")
	ErrBadMoney = errors.New("incorrect amount of money")
)

type Cash struct {
	Money int
}

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

type Wallet struct {
	Money int
}

func (w *Wallet) Spend(n int) bool {
	if w.Money < n {
		return false
	}
	w.Money -= n
	return true
}

func (w Wallet) String() string {
	return fmt.Sprintf("💰 Кошелек (%d %s)", w.Money, Currency)
}
