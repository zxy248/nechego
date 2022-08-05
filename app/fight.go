package app

import (
	"errors"
	"fmt"
	"nechego/model"
	"nechego/service"
	"nechego/statistics"

	tele "gopkg.in/telebot.v3"
)

const profile = Response(`📇 <b>%s %s</b>
<code>%s  %s  %s  %s  %s</code>

💵 Денег в кошельке: %s
💳 На счету в банке: %s

%s

%s`)

// !профиль
func (a *App) handleProfile(c tele.Context) error {
	user := getUser(c)
	strength, err := a.stat.Strength(user)
	if err != nil {
		return respondInternalError(c, err)
	}
	modset, err := a.stat.UserModset(user)
	if err != nil {
		return respondInternalError(c, err)
	}
	return respond(c, profile.Fill(
		formatTitles(modset.Titles()...),
		a.mustMentionUser(user),
		formatEnergy(user.Energy),
		formatElo(user.Elo),
		formatStrength(strength),
		formatMessages(user.Messages),
		formatFood(user.Fishes),
		formatMoney(user.Balance),
		formatMoney(user.Account),
		formatStatus(modset.Descriptions()...),
		formatIcons(modset.Icons()...),
	))
}

const (
	versus               = "⚔️ <b>%s</b> <code>[%.2f]</code> <b><i>vs</i></b> <b>%s</b> <code>[%.2f]</code>"
	fightCollect         = "🏆 %s <code>(%s)</code> выигрывает в поединке и забирает %s"
	fightNoMoney         = "🏆 %s <code>(%s)</code> выигрывает в поединке."
	cannotAttackYourself = UserError("Вы не можете напасть на самого себя.")
)

// !драка
func (a *App) handleFight(c tele.Context) error {
	outcome, err := a.service.Fight(getUser(c), getReplyUser(c))
	if err != nil {
		if errors.Is(err, service.ErrSameUser) {
			return respondUserError(c, cannotAttackYourself)
		}
		if errors.Is(err, service.ErrNotEnoughEnergy) {
			return respondUserError(c, notEnoughEnergy)
		}
		return respondInternalError(c, err)
	}
	return respond(c, a.fightResponse(outcome))
}

func (a *App) fightResponse(o *service.FightOutcome) Response {
	sections := []string{versus}
	args := []any{
		a.mustMentionUser(o.Attacker.User),
		o.Attacker.Strength,
		a.mustMentionUser(o.Defender.User),
		o.Defender.Strength,
		a.mustMentionUser(o.Winner().User),
		formatEloDelta(o.Elo),
	}
	if o.Reward > 0 {
		args = append(args, formatMoney(o.Reward))
		sections = append(sections, fightCollect)
	} else {
		sections = append(sections, fightNoMoney)
	}
	sections = append(sections, string(energyRemaining(o.Attacker.Energy)))
	return Response(joinSections(sections...)).Fill(args...)
}

const (
	topStrong = Response(`🏋️‍♀️ <b>Самые сильные пользователи</b>
%s`)
	topWeak = Response(`🤕 <b>Самые слабые пользователи</b>
%s`)
)

// !топ сильных
func (a *App) handleTopStrong(c tele.Context) error {
	users, err := a.stat.SortedUsers(getGroup(c), a.stat.ByStrengthDesc)
	if err != nil {
		return respondInternalError(c, err)
	}
	n := clampTopNumber(len(users))
	top, err := a.topStrength(users[:n])
	if err != nil {
		return respondInternalError(c, err)
	}
	return respond(c, topStrong.Fill(top))
}

// !топ слабых
func (a *App) handleTopWeak(c tele.Context) error {
	users, err := a.stat.SortedUsers(getGroup(c), a.stat.ByStrengthAsc)
	if err != nil {
		return respondInternalError(c, err)
	}
	n := clampTopNumber(len(users))
	top, err := a.topStrength(users[:n])
	if err != nil {
		return respondInternalError(c, err)
	}
	return respond(c, topWeak.Fill(top))
}

func (a *App) topStrength(u []model.User) (HTML, error) {
	s := []string{}
	for _, uu := range u {
		str, err := a.stat.Strength(uu)
		if err != nil {
			return "", err
		}
		s = append(s, fmt.Sprintf("%s %s", a.mustMentionUser(uu), formatStrength(str)))
	}
	return enumerate(s...), nil
}

// !сила
func (a *App) handleStrength(c tele.Context) error {
	str, err := a.stat.Strength(getUser(c))
	if err != nil {
		return respondInternalError(c, err)
	}
	return respond(c, Response("Ваша сила: %s").Fill(formatStrength(str)))
}

const topRating = Response("🏆 <b>Боевой рейтинг</b>\n%s")

func (a *App) handleTopElo(c tele.Context) error {
	users, err := a.stat.SortedUsers(getGroup(c), statistics.ByEloDesc)
	if err != nil {
		return respondInternalError(c, err)
	}
	users = users[:clampTopNumber(len(users))]
	return respond(c, topRating.Fill(a.topElo(users)))
}

func (a *App) topElo(u []model.User) HTML {
	s := []string{}
	for _, uu := range u {
		s = append(s, fmt.Sprintf("%s %s", a.mustMentionUser(uu), formatElo(uu.Elo)))
	}
	return enumerate(s...)
}
