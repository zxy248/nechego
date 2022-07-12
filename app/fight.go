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
	final, err := a.userStrength(u)
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

func (f fight) sameIDs() bool {
	return f.attacker.ID == f.defender.ID
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
	fightCollect         = "⚔️ *%s* `[%.2f]` _против_ *%s* `[%.2f]`\n\n🏆 Побеждает %s и забирает %s"
	fightNoMoney         = "⚔️ *%s* `[%.2f]` _против_ *%s* `[%.2f]`\n\n🏆 Побеждает %s\\. У проигравшего нечего отнять\\."
	cannotAttackYourself = "Вы не можете напасть на самого себя."
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
	f := fight{attacker, defender}
	if f.sameIDs() {
		return userError(c, cannotAttackYourself)
	}

	ok := a.model.UpdateEnergy(f.attacker.User, -energyDelta, energyCap)
	if !ok {
		return userError(c, notEnoughEnergy)
	}

	win := randInRange(minWinReward, maxWinReward)
	reward, err := a.model.ForceTransferMoney(f.loser().User, f.winner().User, win)
	if err != nil {
		return internalError(c, err)
	}

	template := fightNoMoney
	args := []interface{}{a.mustMentionUser(f.attacker.User),
		f.attacker.actualStrength,
		a.mustMentionUser(f.defender.User),
		f.defender.actualStrength,
		a.mustMentionUser(f.winner().User),
	}
	if reward > 0 {
		template = fightCollect
		args = append(args, formatMoney(reward))
	}
	out := fmt.Sprintf(template, args...)
	out = appendEnergyRemaining(out, f.attacker.Energy-energyDelta)
	return c.Send(out, tele.ModeMarkdownV2)
}

func fightChance() float64 {
	return rand.Float64()*2 - 1
}

const chanceRatio = 0.5

func fightFormula(strength, chance float64) float64 {
	return (strength * (1 - chanceRatio)) + (strength * chance * chanceRatio)
}

// userStrength determines the final strength of a user.
func (a *App) userStrength(u model.User) (float64, error) {
	strength, err := a.actualUserStrength(u)
	if err != nil {
		return 0, err
	}
	return fightFormula(strength, fightChance()), nil
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

// totalMessageCount returns a total message count in the group.
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

const topStrong = "🏋️‍♀️ *Самые сильные пользователи*\n"

// !топ сильных
func (a *App) handleTopStrong(c tele.Context) error {
	users, err := a.strongestUsers(getGroup(c))
	if err != nil {
		return internalError(c, err)
	}
	n := topNumber(len(users))
	strong := users[:n]
	top, err := a.formatTopStrength(strong)
	if err != nil {
		return internalError(c, err)
	}
	return c.Send(topStrong+top, tele.ModeMarkdownV2)
}

const topWeak = "🤕 *Самые слабые пользователи*\n"

// !топ слабых
func (a *App) handleTopWeak(c tele.Context) error {
	users, err := a.strongestUsers(getGroup(c))
	if err != nil {
		return internalError(c, err)
	}
	n := topNumber(len(users))
	weak := []model.User{}
	for i := 0; i < n; i++ {
		weak = append(weak, users[len(users)-1-i])
	}
	top, err := a.formatTopStrength(weak)
	if err != nil {
		return internalError(c, err)
	}
	return c.Send(topWeak+top, tele.ModeMarkdownV2)
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
		var x, y float64
		x, err = a.actualUserStrength(users[i])
		if err != nil {
			return false
		}
		y, err = a.actualUserStrength(users[j])
		if err != nil {
			return false
		}
		return x > y
	})
	return users, err
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
	luck := luckModifier(u)
	if luck != noModifier {
		return append(m, luck), nil
	}
	return m, nil
}

func (a *App) addRichModifier(u model.User, m []*modifier) ([]*modifier, error) {
	rich, err := a.isRich(u)
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

func luckModifier(u model.User) *modifier {
	switch luck := luckLevel(u); {
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
