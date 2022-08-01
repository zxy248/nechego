package statistics

import (
	"math"
	"nechego/model"
)

const EnergyHardLimit = math.MaxInt

func (s *Statistics) HasFullEnergy(u model.User) bool {
	return u.Energy >= s.Settings.EnergyRange.Max()
}

func (s *Statistics) HasMuchEnergy(u model.User) bool {
	return u.Energy > s.Settings.EnergyRange.Max()
}

func (s *Statistics) HasNoEnergy(u model.User) bool {
	return u.Energy == s.Settings.EnergyRange.Min()
}
