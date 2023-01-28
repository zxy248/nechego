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
	if c.Money >= n {
		c.Money -= n
		return true
	}
	return false
}

func (c Cash) String() string {
	return fmt.Sprintf("💵 Наличные (%d %s)", c.Money, Currency)
}

type Wallet struct {
	Money int
}

func (w *Wallet) Spend(n int) bool {
	if w.Money >= n {
		w.Money -= n
		return true
	}
	return false
}

func (w Wallet) String() string {
	return fmt.Sprintf("💰 Кошелек (%d %s)", w.Money, Currency)
}
