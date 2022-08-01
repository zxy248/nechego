package statistics

import (
	"fmt"
	"nechego/model"
	"nechego/modifiers"
	"nechego/numbers"
	"strings"
	"time"
)

func (s *Statistics) UserModset(u model.User) (*modifiers.Set, error) {
	setters := []modifierSetter{
		setAdminModifier,
		s.setEblanModifier,
		s.setEnergyModifier,
		setLuckModifier,
		s.setRichModifier,
		s.setPoorModifier,
		setFisherModifier,
		setDebtorModifier,
	}
	m := modifiers.NewSet()
	for _, set := range setters {
		if err := set(m, u); err != nil {
			return nil, err
		}
	}
	return m, nil
}

type modifierSetter func(*modifiers.Set, model.User) error

func setAdminModifier(m *modifiers.Set, u model.User) error {
	if u.Admin {
		m.Add(modifiers.AdminModifier)
	}
	return nil
}

func (s *Statistics) setEblanModifier(m *modifiers.Set, u model.User) error {
	eblan, err := s.model.DailyEblan(model.Group{GID: u.GID})
	if err != nil {
		return err
	}
	if eblan.ID == u.ID {
		m.Add(modifiers.EblanModifier)
	}
	return nil
}

func (s *Statistics) setEnergyModifier(m *modifiers.Set, u model.User) error {
	m.Add(s.energyModifier(u))
	return nil
}

func (s *Statistics) energyModifier(u model.User) *modifiers.Modifier {
	if s.HasMuchEnergy(u) {
		return modifiers.MuchEnergyModifier
	}
	if s.HasFullEnergy(u) {
		return modifiers.FullEnergyModifier
	}
	if s.HasNoEnergy(u) {
		return modifiers.NoEnergyModifier
	}
	return modifiers.NoModifier
}

func setLuckModifier(m *modifiers.Set, u model.User) error {
	m.Add(luckModifier(u))
	return nil
}

func luckModifier(u model.User) *modifiers.Modifier {
	switch luck := luckLevel(u); {
	case luck < 10:
		return modifiers.TerribleLuckModifier
	case luck < 40:
		return modifiers.BadLuckModifier
	case luck < 70:
		return modifiers.GoodLuckModifier
	case luck < 80:
		return modifiers.ExcellentLuckModifier
	}
	return modifiers.NoModifier
}

func luckLevel(u model.User) byte {
	now := time.Now()
	values := []any{u.UID, u.GID, now.Day(), now.Month(), now.Year()}
	format := strings.Repeat("%v", len(values))
	seed := fmt.Sprintf(format, values...)
	return numbers.LuckyByte([]byte(seed))
}

func (s *Statistics) setRichModifier(m *modifiers.Set, u model.User) error {
	rich, err := s.IsRich(u)
	if err != nil {
		return err
	}
	if rich {
		m.Add(modifiers.RichModifier)
	}
	return nil
}

func (s *Statistics) setPoorModifier(m *modifiers.Set, u model.User) error {
	if s.IsPoor(u) {
		m.Add(modifiers.PoorModifier)
	}
	return nil
}

func setFisherModifier(m *modifiers.Set, u model.User) error {
	if u.Fisher {
		m.Add(modifiers.FisherModifier)
	}
	return nil
}

func setDebtorModifier(m *modifiers.Set, u model.User) error {
	if u.Debtor() {
		m.Add(modifiers.DebtorModifier)
	}
	return nil
}
