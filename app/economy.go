package app

import (
	"errors"
	"fmt"
	"nechego/input"
	"nechego/model"
	"nechego/service"
	"nechego/statistics"
	"strings"

	tele "gopkg.in/telebot.v3"
)

const (
	notEnoughMoney      = UserError("Недостаточно средств.")
	notEnoughMoneyDelta = UserError("Вам не хватает %s")
	specifyAmount       = UserError("Укажите корректное количество средств.")
	incorrectAmount     = UserError("Некорректная сумма.")
)

// !баланс
func (a *App) handleBalance(c tele.Context) error {
	return respond(c, balanceResponse(getUser(c)))
}

func balanceResponse(u model.User) Response {
	s := []any{inTheWallet(u.Balance), onTheAccount(u.Account)}
	if u.Debtor() {
		s = append(s, debtValue(u.Debt))
	}
	return Response(strings.Repeat("%s\n", len(s))).Fill(s...)
}

func inTheWallet(n int) HTML {
	return HTML(fmt.Sprintf("💵 В кошельке: %s", formatMoney(n)))
}

func onTheAccount(n int) HTML {
	return HTML(fmt.Sprintf("💳 На банковском счете: %s", formatMoney(n)))
}

func debtValue(n int) HTML {
	return HTML(fmt.Sprintf("🏦 Кредит: %s", formatMoney(n)))
}

const transfer = Response("Вы перевели %s %s")

// !перевод
func (a *App) handleTransfer(c tele.Context) error {
	sender := getUser(c)
	recipient := getReplyUser(c)
	amount, err := getMessage(c).MoneyArgument()
	if errors.Is(err, input.ErrAllIn) {
		amount = sender.Balance
	} else if err != nil {
		return respondUserError(c, specifyAmount)
	}
	if err := a.service.Transfer(sender, recipient, amount); err != nil {
		var moneyErr service.NotEnoughMoneyError
		if errors.As(err, &moneyErr) {
			return respondUserError(c, notEnoughMoneyDelta.Fill(formatMoney(moneyErr.Delta)))
		}
		return respondInternalError(c, err)
	}
	return respond(c, transfer.Fill(a.mustMentionUser(recipient), formatMoney(amount)))
}

const capital = Response(`💸 Капитал беседы <b>%s</b>: %s


<i>В руках магната %s %s,</i>
<i>или %s от общего количества средств.</i>

<i>В среднем на счету у пользователя: %s</i>`)

func (a *App) handleCapital(c tele.Context) error {
	group := getGroup(c)
	title := c.Chat().Title
	richest, err := a.stat.GreatestUser(group, statistics.ByWealthDesc)
	if err != nil {
		return respondInternalError(c, err)
	}
	total, err := a.stat.GroupBalance(group)
	if err != nil {
		return respondInternalError(c, err)
	}
	average, err := a.stat.AverageBalance(group)
	if err != nil {
		return respondInternalError(c, err)
	}
	percentage := float64(richest.Summary()) / float64(total)
	return respond(c, capital.Fill(
		title,
		formatMoney(total),
		a.mustMentionUser(richest),
		formatMoney(richest.Summary()),
		formatPercentage(percentage),
		formatMoney(int(average))))
}

const (
	topRich = Response(`💵 <b>Самые богатые пользователи</b>
%s`)
	topPoor = Response(`🗑 <b>Самые бедные пользователи</b>
%s`)
)

// !топ богатых
func (a *App) handleTopRich(c tele.Context) error {
	users, err := a.stat.SortedUsers(getGroup(c), statistics.ByWealthDesc)
	if err != nil {
		return respondInternalError(c, err)
	}
	n := clampTopNumber(len(users))
	return respond(c, topRich.Fill(a.topRich(users[:n])))
}

// !топ нищих
func (a *App) handleTopPoor(c tele.Context) error {
	users, err := a.stat.SortedUsers(getGroup(c), statistics.ByWealthAsc)
	if err != nil {
		return respondInternalError(c, err)
	}
	n := clampTopNumber(len(users))
	return respond(c, topPoor.Fill(a.topRich(users[:n])))
}

func (a *App) topRich(u []model.User) HTML {
	s := []string{}
	for _, uu := range u {
		s = append(s, fmt.Sprintf("%s %s", a.mustMentionUser(uu), formatMoney(uu.Summary())))
	}
	return enumerate(s...)
}
