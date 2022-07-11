package app

import (
	"errors"
	"fmt"
	"nechego/model"

	tele "gopkg.in/telebot.v3"
)

const bank = "🏦 *Банк:* на вашем счете %s\n\n" +
	"_Снять средства: `!обнал`\\._\n" +
	"_Пополнить счет: `!депозит`\\._\n\n" +
	"_%s_\n\n" +
	"_Взять кредит: `!кредит`\\._\n" +
	"_Погасить кредит: `!погасить`\\._\n" +
	"_Процентная ставка: %s_\n" +
	"_Кредитный лимит: %s_\n" +
	"_Комиссия за пополнение: %s_\n"

func (a *App) handleBank(c tele.Context) error {
	user := getUser(c)
	return c.Send(fmt.Sprintf(bank,
		formatMoney(user.Account),
		debtStatus(user),
		formatRatio(debtFee),
		formatMoney(user.DebtLimit),
		formatMoney(bankFee)),
		tele.ModeMarkdownV2)
}

func debtStatus(u model.User) string {
	if !u.Debtor() {
		return "У вас нет кредитов\\."
	}
	return fmt.Sprintf("Вы должны банку %s", formatMoney(u.Debt))
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

const (
	withdraw     = "💳 Вы сняли %s со счета\\.\n\n_Теперь в кошельке %s_"
	withdrawDebt = "Вы не можете снимать средства со счета, пока у вас есть непогашенные кредиты.\n"
)

func (a *App) handleWithdraw(c tele.Context) error {
	user := getUser(c)
	if user.Debtor() {
		return userError(c, withdrawDebt)
	}
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

const (
	debtorCannotLoan = "Вы не можете взять средства в долг, пока у вас есть непогашенные кредиты."
	debtTooLow       = "Минимальный кредит — %s"
	limitTooLow      = "Ваш кредитный лимит — %s"
	debtSuccess      = "💳 Вы взяли в кредит %s\n\n_Вам необходимо вернуть %s_"
)

// !долг, !кредит
func (a *App) handleDebt(c tele.Context) error {
	user := getUser(c)
	if user.Debtor() {
		return userError(c, debtorCannotLoan)
	}
	amount, err := moneyArgument(c)
	if amount == 0 || err != nil {
		return err
	}
	if amount < minDebt {
		return userErrorMarkdown(c, fmt.Sprintf(debtTooLow, formatMoney(minDebt)))
	}
	fee := int(float64(amount) * debtFee)
	ok := a.model.Loan(user, amount, fee)
	if !ok {
		return userErrorMarkdown(c, fmt.Sprintf(limitTooLow, formatMoney(user.DebtLimit)))
	}
	return c.Send(fmt.Sprintf(debtSuccess, formatMoney(amount), formatMoney(amount+fee)),
		tele.ModeMarkdownV2)
}

const (
	notDebtor              = "У вас нет непогашенных кредитов."
	notEnoughOnBankAccount = "Недостаточно средств на банковском счете."
	repayFullSuccess       = "💳 Вы погасили кредит."
	repayPartialSuccess    = "💳 Вы погасили %s\n\n_Осталось погасить: %s_"
)

func (a *App) handleRepay(c tele.Context) error {
	user := getUser(c)
	if !user.Debtor() {
		return userError(c, notDebtor)
	}
	amount, err := moneyArgument(c)
	if amount == 0 || err != nil {
		return err
	}
	if user.Debt <= amount {
		amount = user.Debt
	}
	ok := a.model.Repay(user, amount)
	if !ok {
		return userError(c, notEnoughOnBankAccount)
	}
	if amount == user.Debt {
		return c.Send(repayFullSuccess)
	}
	return c.Send(fmt.Sprintf(repayPartialSuccess, formatMoney(amount), formatMoney(user.Debt-amount)),
		tele.ModeMarkdownV2)
}
