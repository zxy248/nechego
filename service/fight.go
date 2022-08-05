package service

import (
	"errors"
	"math/rand"
	"nechego/fight"
	"nechego/model"
	"nechego/numbers"
	"nechego/statistics"
)

var ErrSameUser = errors.New("same user")

type FightOutcome struct {
	*fight.Fight
	Reward int
	Elo    float64
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
	winner := battle.Winner().User
	loser := battle.Loser().User
	reward, err := s.model.ForceTransferMoney(loser, winner,
		fightReward(s.Config.MinReward, s.Config.BaseReward, winner.Balance, loser.Balance))
	if err != nil {
		return nil, err
	}
	elo := numbers.EloDelta(winner.Elo, loser.Elo, numbers.KDefault, numbers.ScoreWin)
	s.model.UpdateElo(winner, elo)
	s.model.UpdateElo(loser, -elo)
	return &FightOutcome{
		Fight:  battle,
		Reward: reward,
		Elo:    elo,
	}, nil
}

const (
	winnerRewardFactor = 1.0
	loserRewardFactor  = 0.125
	sigmaRewardFactor  = 0.5
)

func fightReward(minReward, baseReward, winnerBalance, loserBalance int) int {
	return int(fightRewardHelper(
		float64(minReward),
		float64(baseReward),
		float64(winnerBalance),
		float64(loserBalance),
	))
}

func fightRewardHelper(minReward, baseReward, winnerBalance, loserBalance float64) float64 {
	x := numbers.Max(
		baseReward,
		numbers.Min(
			winnerBalance*winnerRewardFactor,
			loserBalance*loserRewardFactor,
		),
	)
	return numbers.Max(
		minReward,
		rand.NormFloat64()*sigmaRewardFactor*x+x,
	)
}
