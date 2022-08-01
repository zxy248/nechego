package statistics

import (
	"nechego/model"
	"nechego/numbers"
)

type Settings struct {
	EnergyRange   numbers.Interval
	PoorThreshold int
}

type Statistics struct {
	Settings Settings
	model    *model.Model
}

func New(m *model.Model, s Settings) *Statistics {
	return &Statistics{
		Settings: s,
		model:    m,
	}
}
