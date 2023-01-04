package game

import (
	"fmt"
	"time"
)

type Wallet struct {
	Money int
}

func (w Wallet) String() string {
	return fmt.Sprintf("💰 Кошелек (%d ₽)", w.Money)
}

type CreditCard struct {
	Bank    int
	Number  int
	Expires time.Time
	Money   int
}

func (c CreditCard) String() string {
	return fmt.Sprintf("💳 Кредитная карта (%d ₽)", c.Money)
}

type Debt struct {
	CreditorID int
	Money      int
	Percent    int
}

func (d Debt) String() string {
	return fmt.Sprintf("💵 Долг (%d ₽, %d%%)", d.Money, d.Percent)
}
