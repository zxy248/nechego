package app

import (
	"errors"
	"fmt"
	"nechego/input"
	"nechego/model"
	"sort"

	tele "gopkg.in/telebot.v3"
)

// handleBalance responds with the balance of a user.
func (a *App) handleBalance(c tele.Context) error {
	user := getUser(c)
	out := inTheWallet(user.Balance) + onTheAccount(user.Account)
	if user.Debtor() {
		out += debtValue(user.Debt)
	}
	return c.Send(out, tele.ModeMarkdownV2)
}

func inTheWallet(n int) string {
	return fmt.Sprintf("💵 В кошельке: %s\n", formatMoney(n))
}

func onTheAccount(n int) string {
	return fmt.Sprintf("💳 На банковском счете: %s\n", formatMoney(n))
}

func debtValue(n int) string {
	return fmt.Sprintf("🏦 Кредит: %s\n", formatMoney(n))
}

const (
	handleTransferTemplate = "Вы перевели %s %s"
	notEnoughMoney         = "Недостаточно средств."
	specifyAmount          = "Укажите количество средств."
	incorrectAmount        = "Некорректная сумма."
)

// handleTransfer transfers the specified amount of money from one user to another.
func (a *App) handleTransfer(c tele.Context) error {
	sender := getUser(c)
	recipient := getReplyUser(c)
	amount, err := moneyArgument(c)
	if amount == 0 || err != nil {
		return err
	}

	if err := a.model.TransferMoney(sender, recipient, amount); err != nil {
		if errors.Is(err, model.ErrNotEnoughMoney) {
			return userError(c, notEnoughMoney)
		}
		return internalError(c, err)
	}
	out := fmt.Sprintf(handleTransferTemplate, a.mustMentionUser(recipient), formatMoney(amount))
	return c.Send(out, tele.ModeMarkdownV2)
}

// check if int == 0
func moneyArgument(c tele.Context) (int, error) {
	amount, err := getMessage(c).MoneyArgument()
	if err != nil {
		if errors.Is(err, input.ErrSpecifyAmount) {
			return 0, userError(c, specifyAmount)
		}
		if errors.Is(err, input.ErrNotPositive) {
			return 0, userError(c, incorrectAmount)
		}
		return 0, internalError(c, err)
	}
	return amount, nil
}

// isPoor returns true if the user's wealth is less than the maximum win reward.
func isPoor(u model.User) bool {
	return u.Summary() < maxWinReward
}

// isRich returns true if the user is the richest user in the group.
func (a *App) isRich(u model.User) (bool, error) {
	group, err := a.model.GetGroup(model.Group{GID: u.GID})
	if err != nil {
		return false, err
	}
	richest, err := a.richestUser(group)
	if err != nil {
		return false, err
	}
	return richest.ID == u.ID, nil
}

// richestUser returns the richest user in the group.
func (a *App) richestUser(g model.Group) (model.User, error) {
	users, err := a.richestUsers(g)
	if err != nil {
		return model.User{}, nil
	}
	if len(users) < 1 {
		return model.User{}, errors.New("list of users is too short")
	}
	return users[0], nil
}

// richestUsers returns a list of users in the group sorted by wealth.
func (a *App) richestUsers(g model.Group) ([]model.User, error) {
	users, err := a.model.ListUsers(g)
	if err != nil {
		return nil, err
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Summary() > users[j].Summary()
	})
	return users, nil
}

const handleProfileTemplate = `ℹ️ *Профиль %s %v %s*

Денег в кошельке: %s
На счету в банке: %s
Запас энергии: %s
Базовая сила: %s
Написано сообщений: %s
Имеется рыбы: %s

%s
`

// handleProfile sends the profile of the user.
func (a *App) handleProfile(c tele.Context) error {
	user := getUser(c)
	icon := "👤"
	title := "пользователя"

	strength, err := a.actualUserStrength(user)
	if err != nil {
		return internalError(c, err)
	}

	var status string
	modifiers, err := a.userModifiers(user)
	if err != nil {
		return internalError(c, err)
	}
	for _, m := range modifiers {
		switch m {
		case eblanModifier:
			icon, title = "😸", "еблана"
		case adminModifier:
			icon, title = "👑", "администратора"
		case terribleLuckModifier:
			icon = "☠️"
		case excellentLuckModifier:
			icon = "🍀"
		case richModifier:
			icon, title = "🎩", "магната"
		}
		if m != noModifier {
			status += m.description + "\n"
		}
	}
	if status != "" {
		status = fmt.Sprintf("_%s_", markdownEscaper.Replace(status))
	}

	out := fmt.Sprintf(handleProfileTemplate,
		title, a.mustMentionUser(user), icon,
		formatMoney(user.Balance),
		formatMoney(user.Account),
		formatEnergy(user.Energy),
		formatStrength(strength),
		formatMessages(user.Messages),
		formatFishes(user.Fishes),
		status)
	return c.Send(out, tele.ModeMarkdownV2)
}

const topRich = "💵 *Самые богатые пользователи*\n"

// handleTopRich sends a top of the richest users.
func (a *App) handleTopRich(c tele.Context) error {
	users, err := a.richestUsers(getGroup(c))
	if err != nil {
		return internalError(c, err)
	}
	n := topNumber(len(users))
	rich := users[:n]
	top := a.formatTopRich(rich)
	return c.Send(topRich+top, tele.ModeMarkdownV2)
}

const topPoor = "🗑 *Самые бедные пользователи*\n"

// handleTopPoor sends a top of the poorest users.
func (a *App) handleTopPoor(c tele.Context) error {
	users, err := a.richestUsers(getGroup(c))
	if err != nil {
		return internalError(c, err)
	}
	n := topNumber(len(users))
	poor := []model.User{}
	for i := 0; i < n; i++ {
		poor = append(poor, users[len(users)-1-i])
	}
	top := a.formatTopRich(poor)
	return c.Send(topPoor+top, tele.ModeMarkdownV2)
}

const handleCapitalTemplate = "💸 Капитал беседы *%s*: %s\n\n" +
	"_В руках магната %s %s,\nили `%.1f%%` от общего количества средств\\._\n\n" +
	"_В среднем на счету у пользователя: %s_\n"

func (a *App) handleCapital(c tele.Context) error {
	group := getGroup(c)
	title := c.Chat().Title
	richest, err := a.richestUser(group)
	if err != nil {
		return internalError(c, err)
	}
	total, err := a.groupBalance(group)
	if err != nil {
		return internalError(c, err)
	}
	average, err := a.averageBalance(group)
	if err != nil {
		return internalError(c, err)
	}
	percentage := float64(richest.Summary()) / float64(total) * 100
	out := fmt.Sprintf(handleCapitalTemplate,
		title, formatMoney(total),
		a.mustMentionUser(richest), formatMoney(richest.Summary()), percentage,
		formatMoney(int(average)))
	return c.Send(out, tele.ModeMarkdownV2)
}

// groupBalance returns the group's balance.
func (a *App) groupBalance(g model.Group) (int, error) {
	users, err := a.model.ListUsers(g)
	if err != nil {
		return 0, err
	}
	sum := 0
	for _, u := range users {
		sum += u.Summary()
	}
	return sum, nil
}

// averageBalance returns the group's average balance.
func (a *App) averageBalance(g model.Group) (float64, error) {
	users, err := a.model.ListUsers(g)
	if err != nil {
		return 0, err
	}
	if len(users) == 0 {
		return 0, errors.New("list of users is empty")
	}
	sum := 0
	for _, u := range users {
		sum += u.Summary()
	}
	return float64(sum) / float64(len(users)), nil
}
