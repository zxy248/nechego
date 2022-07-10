package app

import (
	"errors"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

const bank = "🏦 *Банк:* на вашем счете %s\n\n" +
	"_Снять средства: `!обнал`\\._\n" +
	"_Пополнить счет: `!депозит`\\._\n" +
	"_Комиссия за пополнение: %s_"

func (a *App) handleBank(c tele.Context) error {
	return c.Send(fmt.Sprintf(bank,
		formatMoney(getUser(c).Account),
		formatMoney(bankFee)),
		tele.ModeMarkdownV2)
}

const deposit = "💳 Вы оплатили налог и положили %s в банк\\.\n\n_Теперь на счету %s_"

func (a *App) handleDeposit(c tele.Context) error {
	user := getUser(c)
	amount, err := moneyArgument(c)
	if amount == 0 || err != nil {
		return err
	}
	amount, err = amountAfterBankFee(amount)
	if err != nil {
		return userError(c, err.Error())
	}
	ok := a.model.Deposit(user, amount, bankFee)
	if !ok {
		return userError(c, notEnoughMoney)
	}
	return c.Send(fmt.Sprintf(deposit, formatMoney(amount), formatMoney(user.Account+amount)),
		tele.ModeMarkdownV2)
}

const withdraw = "💳 Вы сняли %s со счета\\.\n\n_Теперь в кошельке %s_"

func (a *App) handleWithdraw(c tele.Context) error {
	user := getUser(c)
	amount, err := moneyArgument(c)
	if amount == 0 || err != nil {
		return err
	}
	ok := a.model.Withdraw(user, amount, 0)
	if !ok {
		return userError(c, notEnoughMoney)
	}
	return c.Send(fmt.Sprintf(withdraw, formatMoney(amount), formatMoney(user.Balance+amount)),
		tele.ModeMarkdownV2)
}

func amountAfterBankFee(amount int) (int, error) {
	amount = amount - bankFee
	if amount <= 0 {
		return 0, errors.New(notEnoughMoney)
	}
	return amount, nil
}
