package service

import (
	"errors"
)

var ErrNotEnoughEnergy = errors.New("not enough energy")

func (s *Service) RestoreEnergy() {
	s.model.RestoreEnergy(s.Config.EnergyRestoreDelta, s.statistics.Settings.EnergyRange.Max())
}
