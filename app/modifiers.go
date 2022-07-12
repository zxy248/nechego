package app

import (
	"crypto/sha1"
	"fmt"
	"nechego/model"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

type modifier struct {
	icon        string
	multiplier  float64
	description string
}

type modset []*modifier

func newModset() *modset {
	return &modset{}
}

func (ms *modset) present(m *modifier) bool {
	return slices.Contains(*ms, m)
}

func (ms *modset) add(m *modifier) {
	if slices.Contains(*ms, m) {
		return
	}
	if m == noModifier {
		return
	}
	*ms = append(*ms, m)
}

func (ms *modset) list() []*modifier {
	return *ms
}

func (ms *modset) sum() float64 {
	sum := float64(0)
	for _, m := range ms.list() {
		sum += m.multiplier
	}
	return sum
}

var (
	noModifier            = &modifier{"", +0.00, ""}
	adminModifier         = &modifier{"👑", +0.20, "Вы ощущаете власть над остальными."}
	eblanModifier         = &modifier{"😸", -0.20, "Вы чувствуете себя оскорбленным."}
	muchEnergyModifier    = &modifier{"🍥", +0.20, "Вы хорошо поели."}
	fullEnergyModifier    = &modifier{"⚡️", +0.10, "Вы полны сил."}
	noEnergyModifier      = &modifier{"😣", -0.25, "Вы чувствуете себя уставшим."}
	terribleLuckModifier  = &modifier{"☠️", -0.50, "Вас преследуют неудачи."}
	badLuckModifier       = &modifier{"", -0.10, "Вам не везет."}
	goodLuckModifier      = &modifier{"🤞", +0.10, "Вам везет."}
	excellentLuckModifier = &modifier{"🍀", +0.30, "Сегодня ваш день."}
	richModifier          = &modifier{"🎩", +0.05, "Вы богаты."}
	poorModifier          = &modifier{"", -0.05, "Вы бедны."}
	fisherModifier        = &modifier{"🎣", +0.05, "Вы можете рыбачить."}
	debtorModifier        = &modifier{"💳", -0.25, "У вас есть кредит."}
)

// userModset returns the user's modset
func (a *App) userModset(u model.User) (*modset, error) {
	setters := []modsetter{
		setAdminModifier,
		a.setEblanModifier,
		setEnergyModifier,
		setLuckModifier,
		a.setRichModifier,
		setPoorModifier,
		setFisherModifier,
		setDebtorModifier,
	}
	ms := newModset()
	for _, set := range setters {
		if err := set(ms, u); err != nil {
			return nil, err
		}
	}
	return ms, nil
}

type modsetter func(*modset, model.User) error

func setAdminModifier(ms *modset, u model.User) error {
	if u.Admin {
		ms.add(adminModifier)
	}
	return nil
}

func (a *App) setEblanModifier(ms *modset, u model.User) error {
	group, err := a.model.GetGroup(model.Group{GID: u.GID})
	if err != nil {
		return err
	}
	eblan, err := a.model.GetDailyEblan(group)
	if err != nil {
		return err
	}
	if eblan.ID == u.ID {
		ms.add(eblanModifier)
	}
	return nil
}

func setEnergyModifier(ms *modset, u model.User) error {
	ms.add(energyModifier(u))
	return nil
}

// energyModifier returns the user's energy modifier.
// If there is no modifier, returns noModifier, nil.
func energyModifier(u model.User) *modifier {
	if hasMuchEnergy(u) {
		return muchEnergyModifier
	}
	if hasFullEnergy(u) {
		return fullEnergyModifier
	}
	if hasNoEnergy(u) {
		return noEnergyModifier
	}
	return noModifier
}

func setLuckModifier(ms *modset, u model.User) error {
	ms.add(luckModifier(u))
	return nil
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

func luckLevel(u model.User) byte {
	now := time.Now()
	values := []any{u.UID, u.GID, now.Day(), now.Month(), now.Year()}
	template := strings.Repeat("%v", len(values))
	seed := fmt.Sprintf(template, values...)
	return sha1.Sum([]byte(seed))[0]
}

func (a *App) setRichModifier(ms *modset, u model.User) error {
	rich, err := a.isRich(u)
	if err != nil {
		return err
	}
	if rich {
		ms.add(richModifier)
	}
	return nil
}

func setPoorModifier(ms *modset, u model.User) error {
	if isPoor(u) {
		ms.add(poorModifier)
	}
	return nil
}

func setFisherModifier(ms *modset, u model.User) error {
	if u.Fisher {
		ms.add(fisherModifier)
	}
	return nil
}

func setDebtorModifier(ms *modset, u model.User) error {
	if u.Debtor() {
		ms.add(debtorModifier)
	}
	return nil
}
