package service

import (
	"errors"
	"nechego/fight"
	"nechego/model"
	"nechego/numbers"
	"nechego/statistics"
)

var ErrSameUser = errors.New("same user")

type FightOutcome struct {
	*fight.Fight
	Reward int
}

func (s *Service) Fight(attacker, defender model.User) (*FightOutcome, error) {
	battle, err := fight.New(attacker, defender, s.Config.FightSettings)
	if err != nil {
		return nil, ErrSameUser
	}
	enoughEnergy := s.model.UpdateEnergy(attacker, -s.Config.FightEnergyDrain, statistics.EnergyHardLimit)
	if !enoughEnergy {
		return nil, ErrNotEnoughEnergy
	}
	battle.Attacker.Energy -= s.Config.FightEnergyDrain
	reward, err := s.model.ForceTransferMoney(
		battle.Loser().User,
		battle.Winner().User,
		numbers.InRange(
			s.Config.WinReward.Min(),
			s.Config.WinReward.Max()))
	if err != nil {
		return nil, err
	}
	return &FightOutcome{
		Fight:  battle,
		Reward: reward,
	}, nil
}
