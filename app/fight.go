package app

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"nechego/model"
	"sort"
	"time"

	tele "gopkg.in/telebot.v3"
)

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
	fightersTemplate          = "⚔️ Нападает %s, сила в бою `%.2f [%.2f]`\n🛡 Защищается %s, сила в бою `%.2f [%.2f]`\n\n"
	winnerTemplate            = "🏆 %s выходит победителем и забирает %s\n\n"
	poorWinnerTemplate        = "🏆 %s выходит победителем и забирает из последних запасов проигравшего %s\n\n"
	energyRemainingTemplate   = "Энергии осталось: %s"
	handleFightTemplate       = fightersTemplate + winnerTemplate + energyRemainingTemplate
	handleFightPoorTemplate   = fightersTemplate + poorWinnerTemplate + energyRemainingTemplate
	displayStrengthMultiplier = 10
	cannotAttackYourself      = "Вы не можете напасть на самого себя."
	notEnoughEnergy           = "Недостаточно энергии."
)

// handleFight conducts a fight between two users.
func (a *App) handleFight(c tele.Context) error {
	attacker, err := a.makeFighter(getUser(c))
	if err != nil {
		return internalError(c, err)
	}
	defender, err := a.makeFighter(getReplyUser(c))
	if err != nil {
		return internalError(c, err)
	}
	if attacker.ID == defender.ID {
		return userError(c, cannotAttackYourself)
	}
	f := fight{attacker, defender}

	ok := a.model.UpdateEnergy(f.attacker.User, -energyDelta, energyCap)
	if !ok {
		return userError(c, notEnoughEnergy)
	}

	win := randInRange(minWinReward, maxWinReward)
	reward, err := a.model.ForceTransferMoney(f.loser().User, f.winner().User, win)
	if err != nil {
		return internalError(c, err)
	}

	var template string
	if reward == 0 {
		reward = randInRange(minWinReward, maxPoorWinReward)
		a.model.UpdateMoney(f.winner().User, reward)
		template = handleFightPoorTemplate
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
		formatMoney(reward),
		formatEnergy(f.attacker.Energy-energyDelta))
	return c.Send(out, tele.ModeMarkdownV2)
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

const handleTopStrength = "🏋️‍♀️ *Самые сильные пользователи*\n%s"

// !топ силы
func (a *App) handleTopStrength(c tele.Context) error {
	users, err := a.strongestUsers(getGroup(c))
	if err != nil {
		return internalError(c, err)
	}
	n := maxTopNumber
	if len(users) < maxTopNumber {
		n = len(users)
	}
	top, err := a.formatTopStrength(users[:n])
	if err != nil {
		return internalError(c, err)
	}
	return c.Send(fmt.Sprintf(handleTopStrength, top), tele.ModeMarkdownV2)
}

func (a *App) strongestUsers(g model.Group) ([]model.User, error) {
	users, err := a.model.ListUsers(g)
	if err != nil {
		return nil, err
	}
	sort.Slice(users, func(i, j int) bool {
		if err != nil {
			return false
		}
		var is, js float64
		is, err = a.actualUserStrength(users[i])
		if err != nil {
			return false
		}
		js, err := a.actualUserStrength(users[j])
		if err != nil {
			return false
		}
		return is > js
	})
	return users, err
}

// !стамина, !энергия
func (a *App) handleEnergy(c tele.Context) error {
	return c.Send(fmt.Sprintf("Осталось энергии: %s", formatEnergy(getUser(c).Energy)), tele.ModeMarkdownV2)
}

// !сила
func (a *App) handleStrength(c tele.Context) error {
	strength, err := a.actualUserStrength(getUser(c))
	if err != nil {
		return internalError(c, err)
	}
	return c.Send(fmt.Sprintf("Ваша сила: %s", formatStrength(strength)), tele.ModeMarkdownV2)
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
	if hasFullEnergy(u) {
		return fullEnergyModifier, nil
	}
	if hasNoEnergy(u) {
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

func (a *App) addFisherModifier(u model.User, m []*modifier) ([]*modifier, error) {
	if u.Fisher {
		return append(m, fisherModifier), nil
	}
	return m, nil
}

func (a *App) addDebtorModifier(u model.User, m []*modifier) ([]*modifier, error) {
	if u.Debtor() {
		return append(m, debtorModifier), nil
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
	fisherModifier        = &modifier{+0.05, "Вы можете рыбачить."}
	debtorModifier        = &modifier{-0.25, "У вас есть кредит."}
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
		a.addFisherModifier,
		a.addDebtorModifier,
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
