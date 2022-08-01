package app

import (
	"errors"
	"fmt"
	"nechego/input"
	"nechego/model"
	"nechego/service"

	tele "gopkg.in/telebot.v3"
)

const bank = Response(`🏦 <b>Банк:</b> на вашем счете %s

<i>Снять средства: <code>!обнал</code></i>
<i>Пополнить счет: <code>!депозит</code></i>
<i>Комиссия на пополнение: %s</i>

<i>%s</i>

<i>Взять кредит: <code>!кредит</code></i>
<i>Погасить кредит: <code>!погасить</code></i>
<i>Процентная ставка: %s</i>
<i>Кредитный лимит: %s</i>`)

// !банк
func (a *App) handleBank(c tele.Context) error {
	user := getUser(c)
	return respond(c, bank.Fill(
		formatMoney(user.Account),
		formatMoney(a.service.Config.DepositFee),
		debtStatus(user),
		formatPercentage(a.service.Config.DebtPercentage),
		formatMoney(user.DebtLimit)))
}

func debtStatus(u model.User) HTML {
	if u.Debtor() {
		return "У вас нет кредитов."
	}
	return HTML(fmt.Sprintf("Вы должны банку %s", formatMoney(u.Debt)))
}

const deposit = Response(`💳 Вы оплатили комиссию и положили %s в банк.

<i>Теперь на счету %s</i>`)

// !депозит
func (a *App) handleDeposit(c tele.Context) error {
	user := getUser(c)
	amount, err := getMessage(c).MoneyArgument()
	if errors.Is(err, input.ErrAllIn) {
		amount = user.Balance
	} else if err != nil {
		return respondUserError(c, specifyAmount)
	}
	transfered, err := a.service.Deposit(user, amount)
	if err != nil {
		if errors.Is(err, service.ErrIncorrectAmount) {
			return respondUserError(c, incorrectAmount)
		}
		var moneyErr service.NotEnoughMoneyError
		if errors.As(err, &moneyErr) {
			return respondUserError(c, notEnoughMoneyDelta.Fill(formatMoney(moneyErr.Delta)))
		}
		return respondInternalError(c, err)
	}
	return respond(c, deposit.Fill(formatMoney(transfered), formatMoney(user.Account+transfered)))
}

const withdraw = Response(`💳 Вы оплатили комиссию и сняли %s со счета.

<i>Теперь в кошельке %s</i>`)

// !обнал
func (a *App) handleWithdraw(c tele.Context) error {
	user := getUser(c)
	amount, err := getMessage(c).MoneyArgument()
	if errors.Is(err, input.ErrAllIn) {
		amount = user.Account
	} else if err != nil {
		return respondUserError(c, specifyAmount)
	}
	transfered, err := a.service.Withdraw(user, amount)
	if err != nil {
		if errors.Is(err, service.ErrIncorrectAmount) {
			return respondUserError(c, incorrectAmount)
		}
		var moneyErr service.NotEnoughMoneyError
		if errors.As(err, &moneyErr) {
			return respondUserError(c, notEnoughMoneyDelta.Fill(formatMoney(moneyErr.Delta)))
		}
		return respondInternalError(c, err)
	}
	return respond(c, withdraw.Fill(formatMoney(transfered), formatMoney(user.Balance+transfered)))
}

const (
	minDebt     = UserError("Минимальный кредит — %s")
	debtLimit   = UserError("Ваш кредитный лимит — %s")
	debtSuccess = Response(`💳 Вы взяли в кредит %s

<i>Вам необходимо вернуть %s</i>`)
)

// !кредит
func (a *App) handleDebt(c tele.Context) error {
	user := getUser(c)
	amount, err := getMessage(c).MoneyArgument()
	if errors.Is(err, input.ErrAllIn) {
		amount = user.DebtLimit
	} else if err != nil {
		return respondUserError(c, specifyAmount)
	}
	debt, err := a.service.Debt(user, amount)
	if err != nil {
		if errors.Is(err, service.ErrDebtLimit) {
			return respondUserError(c, debtLimit.Fill(formatMoney(user.DebtLimit)))
		}
		if errors.Is(err, service.ErrMinDebt) {
			return respondUserError(c, minDebt.Fill(formatMoney(a.service.Config.MinDebt)))
		}
		return respondInternalError(c, err)
	}
	return respond(c, debtSuccess.Fill(formatMoney(amount), formatMoney(debt)))
}

const (
	repayTotalSuccess   = Response("💳 Вы погасили кредит.")
	repayPartialSuccess = Response(`💳 Вы погасили %s

<i>Осталось погасить: %s</i>`)
)

// !погасить
func (a *App) handleRepay(c tele.Context) error {
	user := getUser(c)
	amount, err := getMessage(c).MoneyArgument()
	if errors.Is(err, input.ErrAllIn) {
		amount = user.Account
	} else if err != nil {
		return respondUserError(c, specifyAmount)
	}
	debt, err := a.service.Repay(user, amount)
	if err != nil {
		var moneyErr service.NotEnoughMoneyError
		if errors.As(err, &moneyErr) {
			return respondUserError(c, notEnoughMoneyDelta.Fill(formatMoney(moneyErr.Delta)))
		}
		return respondInternalError(c, err)
	}
	if debt > 0 {
		return respond(c, repayPartialSuccess.Fill(formatMoney(amount), formatMoney(debt)))
	}
	return respond(c, repayTotalSuccess)
}
