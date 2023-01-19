package money

import (
	"errors"
	"fmt"
	"time"
)

const Symbol = "₴"

var (
	ErrNoMoney  = errors.New("insufficient money")
	ErrBadMoney = errors.New("incorrect amount of money")
)

type Cash struct {
	Money int
}

func (c Cash) String() string {
	return fmt.Sprintf("💵 Наличные (%d %s)", c.Money, Symbol)
}

type Wallet struct {
	Money int
}

func (w Wallet) String() string {
	return fmt.Sprintf("💰 Кошелек (%d %s)", w.Money, Symbol)
}

type CreditCard struct {
	Bank    int
	Number  int
	Expires time.Time
	Money   int
}

func (c CreditCard) String() string {
	return fmt.Sprintf("💳 Кредитная карта (%d %s)", c.Money, Symbol)
}

type Debt struct {
	CreditorID int
	Money      int
	Percent    int
}

func (d Debt) String() string {
	return fmt.Sprintf("💵 Долг (%d %s, %d%%)", d.Money, Symbol, d.Percent)
}
