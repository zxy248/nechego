package service

import (
	"nechego/dice"
	"nechego/fight"
	"nechego/model"
	"nechego/numbers"
	"nechego/statistics"
)

type Config struct {
	// fishing
	EatEnergyRestore   numbers.Interval
	FishingRodPrice    int
	FishingEnergyDrain int

	// fighting
	FightSettings    fight.Settings
	FightEnergyDrain int
	MinReward        int
	BaseReward       int

	// parliament
	ParliamentMembers  int
	ParliamentMajority int

	// economy
	MaxDeposits    int
	DepositFee     int
	WithdrawFee    int
	MinDebt        int
	DebtPercentage float64
	InitialBalance int
	InitialElo     float64

	// dice
	MinBet int

	// energy
	EnergyRestoreDelta int

	// pets
	PetPrice int
}

type Service struct {
	Config     Config
	Casino     *dice.Casino
	statistics *statistics.Statistics
	model      *model.Model
}

func New(m *model.Model, s *statistics.Statistics, d *dice.Casino, c Config) *Service {
	return &Service{
		Config:     c,
		Casino:     d,
		statistics: s,
		model:      m,
	}
}
