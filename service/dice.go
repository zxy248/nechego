package service

import (
	"errors"
	"nechego/dice"
	"nechego/model"
)

var (
	ErrBetTooLow = errors.New("bet too low")
)

func (s *Service) Dice(g model.Group, u model.User, bet int, a dice.Actions) error {
	if s.statistics.HasNoEnergy(u) {
		return ErrNotEnoughEnergy
	}
	if bet < s.Config.MinBet {
		return ErrBetTooLow
	}
	enoughMoney := s.model.UpdateMoney(u, -bet)
	if !enoughMoney {
		return ErrNotEnoughMoney
	}
	return s.Casino.Play(g, u, bet, a)
}

func (s *Service) Roll(g model.Group, u model.User, roll int) (dice.Result, error) {
	result, err := s.Casino.Roll(g, u, roll)
	if err != nil {
		return result, err
	}
	switch result.Outcome {
	case dice.Win:
		s.model.UpdateMoney(result.User, result.Bet*2)
	case dice.Draw:
		s.model.UpdateMoney(result.User, result.Bet)
	case dice.Lose:
	}
	return result, nil
}
