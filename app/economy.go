package app

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
	"nechego/input"
	"nechego/model"
	"sort"
	"time"

	tele "gopkg.in/telebot.v3"
)

const handleBalanceTemplate = "Ваш баланс: `%s 💰`"

// handleBalance responds with the balance of a user.
func (a *App) handleBalance(c tele.Context) error {
	return c.Send(fmt.Sprintf(handleBalanceTemplate,
		formatAmount(getUser(c).Balance)),
		tele.ModeMarkdownV2)
}

const handleTransferTemplate = "Вы перевели %s `%s 💰`"

// handleTransfer transfers the specified amount of money from one user to another.
func (a *App) handleTransfer(c tele.Context) error {
	arg, err := getMessage(c).Dynamic()
	if err != nil {
		if errors.Is(err, input.ErrSpecifyAmount) {
			return c.Send(makeError("Укажите количество средств"))
		}
		if errors.Is(err, input.ErrNotPositive) {
			return c.Send(makeError("Некорректная сумма"))
		}
		return err
	}
	amount := arg.(int)

	recipient := getReplyUser(c)
	if err := a.model.TransferMoney(getUser(c), recipient, amount); err != nil {
		if errors.Is(err, model.ErrNotEnoughMoney) {
			return c.Send(makeError("Недостаточно средств"))
		}
		return err
	}
	out := fmt.Sprintf(handleTransferTemplate, a.mustMentionUser(recipient), formatAmount(amount))
	return c.Send(out, tele.ModeMarkdownV2)
}

type fighter struct {
	model.User
	finalStrength  float64
	actualStrength float64
}

func (a *App) makeFighter(u model.User) (fighter, error) {
	final, _, err := a.userStrength(u)
	if err != nil {
		return fighter{}, err
	}
	actual, err := a.actualUserStrength(u)
	if err != nil {
		return fighter{}, err
	}
	return fighter{u, final, actual}, nil
}

type fight struct {
	attacker fighter
	defender fighter
}

func (f fight) winner() fighter {
	if f.attacker.finalStrength > f.defender.finalStrength {
		return f.attacker
	}
	return f.defender
}

func (f fight) loser() fighter {
	if f.attacker.finalStrength <= f.defender.finalStrength {
		return f.attacker
	}
	return f.defender
}

const (
	fightersTemplate = "⚔️ Нападает %s, сила в бою `%.1f [%.1f]`\n" +
		"🛡 Защищается %s, сила в бою `%.1f [%.1f]`\n\n"
	winnerTemplate          = "🏆 %s выходит победителем и забирает `%s 💰`\n\n"
	poorWinnerTemplate      = "🏆 %s выходит победителем и забирает из последних запасов проигравшего `%s 💰`\n\n"
	energyRemainingTemplate = "Энергии осталось: `%v ⚡️`"
	handleFightTemplate     = fightersTemplate + winnerTemplate + energyRemainingTemplate
	handleFightPoorTemplate = fightersTemplate + poorWinnerTemplate + energyRemainingTemplate

	minWinReward              = 1
	maxWinReward              = 10
	maxPoorWinReward          = 3
	displayStrengthMultiplier = 10
)

// handleFight conducts a fight between two users.
func (a *App) handleFight(c tele.Context) error {
	attacker, err := a.makeFighter(getUser(c))
	if err != nil {
		return err
	}
	defender, err := a.makeFighter(getReplyUser(c))
	if err != nil {
		return err
	}
	if attacker.ID == defender.ID {
		return c.Send(makeError("Вы не можете напасть на самого себя"))
	}
	f := fight{attacker, defender}

	ok := a.model.UpdateEnergy(f.attacker.User, -energyDelta, energyCap)
	if !ok {
		return c.Send(makeError("Недостаточно энергии"))
	}

	win := randInRange(minWinReward, maxWinReward)
	transfer, err := a.forceTransferMoney(f.loser().User, f.winner().User, win)
	if err != nil {
		return err
	}

	var template string
	if transfer == 0 {
		reward := randInRange(minWinReward, maxPoorWinReward)
		template = handleFightPoorTemplate
		a.model.UpdateMoney(f.winner().User, reward)
	} else {
		template = handleFightTemplate
	}
	out := fmt.Sprintf(template,
		a.mustMentionUser(f.attacker.User),
		displayStrengthMultiplier*f.attacker.finalStrength,
		f.attacker.actualStrength,
		a.mustMentionUser(f.defender.User),
		displayStrengthMultiplier*f.defender.finalStrength,
		f.defender.actualStrength,
		a.mustMentionUser(f.winner().User),
		formatAmount(transfer),
		f.attacker.Energy-energyDelta)
	return c.Send(out, tele.ModeMarkdownV2)
}

// TODO: !сила

// forceTransferMoney transfers the specified amount of money from one user to another.
// If the sender has not enough money, transfers all the sender's money to the recipient.
func (a *App) forceTransferMoney(sender, recipient model.User, amount int) (int, error) {
	actual := sender.Balance
	if actual < amount {
		return actual, a.model.TransferMoney(sender, recipient, actual)
	}
	return amount, a.model.TransferMoney(sender, recipient, amount)
}

const chanceRatio = 0.5

// userStrength determines the final strength of a user.
func (a *App) userStrength(u model.User) (value float64, chance float64, err error) {
	chance = rand.Float64()*2 - 1
	strength, err := a.actualUserStrength(u)
	if err != nil {
		return 0, 0, err
	}
	result := (strength * (1 - chanceRatio)) + (strength * chance * chanceRatio)
	a.SugarLog().Debugf("(%.1f * (1 - %.1f)) + (%.1f * %.1f * %.1f) = %.1f",
		strength, chanceRatio,
		strength, chance, chanceRatio, result)
	return result, chance, nil
}

const baseStrength = 1

// actualUserStrength determines the user's stength before randomization.
func (a *App) actualUserStrength(u model.User) (float64, error) {
	mcc, err := a.messageCountCoefficient(u)
	if err != nil {
		return 0, err
	}
	mul, err := a.strengthMultiplier(u)
	if err != nil {
		return 0, err
	}
	strength := (baseStrength + mcc) * mul
	return strength, nil
}

// messageCountCoefficient is a quotient of the user's message count and the total message count.
func (a *App) messageCountCoefficient(u model.User) (float64, error) {
	user := u.Messages
	group, err := a.model.GetGroup(model.Group{GID: u.GID})
	if err != nil {
		return 0, err
	}
	total, err := a.totalMessageCount(group)
	if err != nil {
		return 0, err
	}
	return float64(1+user) / float64(1+total), nil
}

// totalMessageCount returns the number of messages sent in the specified interval.
func (a *App) totalMessageCount(g model.Group) (int, error) {
	users, err := a.model.ListUsers(g)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, u := range users {
		total += u.Messages
	}
	return total / len(users), nil
}

// strengthMultiplier returns the strength multiplier value.
func (a *App) strengthMultiplier(u model.User) (float64, error) {
	multiplier := float64(1)
	modifiers, err := a.userModifiers(u)
	if err != nil {
		return 0, err
	}
	for _, m := range modifiers {
		multiplier += m.multiplier
	}
	return multiplier, nil
}

type modifier struct {
	multiplier  float64
	description string
}

type modifierAdder func(u model.User, m []*modifier) ([]*modifier, error)

func (a *App) addAdminModifier(u model.User, m []*modifier) ([]*modifier, error) {
	if u.Admin {
		return append(m, adminModifier), nil
	}
	return m, nil
}

func (a *App) addEblanModifier(u model.User, m []*modifier) ([]*modifier, error) {
	group, err := a.model.GetGroup(model.Group{GID: u.GID})
	if err != nil {
		return nil, err
	}
	eblan, err := a.model.GetDailyEblan(group)
	if err != nil {
		return nil, err
	}
	if eblan.ID == u.ID {
		return append(m, eblanModifier), nil
	}
	return m, nil
}

func (a *App) addEnergyModifier(u model.User, m []*modifier) ([]*modifier, error) {
	energy, err := a.energyModifier(u)
	if err != nil {
		return nil, err
	}
	if energy != noModifier {
		return append(m, energy), nil
	}
	return m, nil
}

// energyModifier returns the user's energy modifier.
// If there is no modifier, returns noModifier, nil.
func (a *App) energyModifier(u model.User) (*modifier, error) {
	if u.Energy == energyCap {
		return fullEnergyModifier, nil
	}
	if u.Energy == 0 {
		return noEnergyModifier, nil
	}
	return noModifier, nil
}

func (a *App) addLuckModifier(u model.User, m []*modifier) ([]*modifier, error) {
	luck := luckModifier(luckLevel(u))
	if luck != noModifier {
		return append(m, luck), nil
	}
	return m, nil
}

func (a *App) addRichModifier(u model.User, m []*modifier) ([]*modifier, error) {
	rich, err := a.isRichest(u)
	if err != nil {
		return nil, err
	}
	if rich {
		return append(m, richModifier), nil
	}
	return m, nil
}

func (a *App) addPoorModifier(u model.User, m []*modifier) ([]*modifier, error) {
	if isPoor(u) {
		return append(m, poorModifier), nil
	}
	return m, nil
}

var (
	noModifier            = &modifier{+0.00, ""}
	adminModifier         = &modifier{+0.20, "Вы ощущаете власть над остальными."}
	eblanModifier         = &modifier{-0.20, "Вы чувствуете себя оскорбленным."}
	fullEnergyModifier    = &modifier{+0.10, "Вы полны сил."}
	noEnergyModifier      = &modifier{-0.25, "Вы чувствуете себя уставшим."}
	terribleLuckModifier  = &modifier{-0.50, "Вас преследуют неудачи."}
	badLuckModifier       = &modifier{-0.10, "Вам не везет."}
	goodLuckModifier      = &modifier{+0.10, "Вам везет."}
	excellentLuckModifier = &modifier{+0.30, "Сегодня ваш день."}
	richModifier          = &modifier{+0.05, "Вы богаты."}
	poorModifier          = &modifier{-0.05, "Вы бедны."}
)

// userModifiers returns the user's modifiers.
func (a *App) userModifiers(u model.User) ([]*modifier, error) {
	adders := []modifierAdder{
		a.addAdminModifier,
		a.addEblanModifier,
		a.addEnergyModifier,
		a.addLuckModifier,
		a.addRichModifier,
		a.addPoorModifier,
	}
	var modifiers []*modifier
	var err error
	for _, add := range adders {
		modifiers, err = add(u, modifiers)
		if err != nil {
			return nil, err
		}
	}
	return modifiers, nil
}

// formatAmount formats the specified amount of money.
func formatAmount(n int) string {
	switch p0 := n % 10; {
	case n >= 10 && n <= 20:
		return fmt.Sprintf("%v монет", n)
	case p0 == 1:
		return fmt.Sprintf("%v монета", n)
	case p0 >= 2 && p0 <= 4:
		return fmt.Sprintf("%v монеты", n)
	default:
		return fmt.Sprintf("%v монет", n)
	}
}

func luckLevel(u model.User) byte {
	now := time.Now()
	seed := fmt.Sprintf("%v%v%v%v%v", u.UID, u.GID, now.Day(), now.Month(), now.Year())
	data := sha1.Sum([]byte(seed))
	return data[0]
}

func luckModifier(luck byte) *modifier {
	switch {
	case luck <= 10:
		return terribleLuckModifier
	case luck <= 40:
		return badLuckModifier
	case luck <= 70:
		return goodLuckModifier
	case luck <= 80:
		return excellentLuckModifier
	}
	return noModifier
}

func isPoor(u model.User) bool {
	return u.Balance < maxWinReward
}

// isRichest returns true if the user is the richest user in the group.
func (a *App) isRichest(u model.User) (bool, error) {
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
		return model.User{}, errors.New("the list of users is too short")
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
		return users[i].Balance > users[j].Balance
	})
	return users, nil
}

// poorestUsers returns a list of users in the group sorted by wealth.
func (a *App) poorestUsers(g model.Group) ([]model.User, error) {
	users, err := a.model.ListUsers(g)
	if err != nil {
		return nil, err
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Balance < users[j].Balance
	})
	return users, nil
}

// TODO: !стамина, !энергия
func handleEnergy(c tele.Context) error {
	return c.Send("У вас %v энергии ⚡️", getUser(c).Energy)
}

const handleProfileTemplate = `ℹ️ Профиль %s %v %s

Баланс на счете: ` + "`" + `%s 💰` + "`" + `
Запас энергии: ` + "`" + `%d ⚡️` + "`" + `
Базовая сила: ` + "`" + `%.2f 💪` + "`" + `
Написано сообщений: ` + "`" + `%d ✍️` + "`" + `

%s
`

// handleProfile sends the profile of the user.
func (a *App) handleProfile(c tele.Context) error {
	user := getUser(c)
	icon := "👤"
	title := "пользователя"

	strength, err := a.actualUserStrength(user)
	if err != nil {
		return err
	}

	var status string
	modifiers, err := a.userModifiers(user)
	if err != nil {
		return err
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
		formatAmount(user.Balance),
		user.Energy,
		strength,
		user.Messages,
		status)
	return c.Send(out, tele.ModeMarkdownV2)
}

const handleTopRichTemplate = "💰 Самые богатые пользователи:%s\n"

// handleTopRich sends a top of the richest users.
func (a *App) handleTopRich(c tele.Context) error {
	users, err := a.richestUsers(getGroup(c))
	if err != nil {
		return err
	}
	n := maxTopNumber
	if len(users) < maxTopNumber {
		n = len(users)
	}
	return c.Send(fmt.Sprintf(handleTopRichTemplate,
		a.formatRichTop(users[:n])), tele.ModeMarkdownV2)
}

const handleTopPoorTemplate = "🗑 Самые бедные пользователи:%s\n"

// handleTopPoor sends a top of the poorest users.
func (a *App) handleTopPoor(c tele.Context) error {
	users, err := a.poorestUsers(getGroup(c))
	if err != nil {
		return err
	}
	n := maxTopNumber
	if len(users) < maxTopNumber {
		n = len(users)
	}
	return c.Send(fmt.Sprintf(handleTopPoorTemplate,
		a.formatRichTop(users[:n])), tele.ModeMarkdownV2)
}

func (a *App) formatRichTop(users []model.User) string {
	var top string
	for i := 0; i < len(users); i++ {
		top += fmt.Sprintf("%d\\. %s, `%s`\n",
			i+1, a.mustMentionUser(users[i]), formatAmount(users[i].Balance))
	}
	return top
}

// TODO: handleTopStrength sends a top of the strongest users.
func handleTopStrength(c tele.Context) error {
	return nil
}

const handleCapitalTemplate = "💸 Капитал беседы *%s*: `%s 💰`\n\n" +
	"_В руках магната %s `%s 💰`,\nили `%.1f%%` от общего количества средств\\._\n\n" +
	"_В среднем на счету у пользователя: `%s 💰`_\n"

func (a *App) handleCapital(c tele.Context) error {
	group := getGroup(c)
	title := c.Chat().Title
	richest, err := a.richestUser(group)
	if err != nil {
		return err
	}
	balance, err := a.groupBalance(group)
	if err != nil {
		return err
	}
	avg, err := a.averageBalance(group)
	if err != nil {
		return err
	}
	percentage := float64(richest.Balance) / float64(balance) * 100
	out := fmt.Sprintf(handleCapitalTemplate,
		title, formatAmount(balance),
		a.mustMentionUser(richest), formatAmount(richest.Balance), percentage,
		formatAmount(int(avg)))
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
		sum += u.Balance
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
		return 0, errors.New("the list of users is empty")
	}
	sum := 0
	for _, u := range users {
		sum += u.Balance
	}
	return float64(sum) / float64(len(users)), nil
}
