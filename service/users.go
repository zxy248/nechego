package service

import (
	"errors"
	"nechego/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func (s *Service) Who(g model.Group, name string) (model.User, error) {
	return s.model.RandomUser(g)
}

func (s *Service) List(g model.Group, n int) ([]model.User, error) {
	return s.model.RandomUsers(g, n)
}

func (s *Service) User(u model.User) (model.User, error) {
	uu, err := s.model.GetUser(u)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			u.Energy = s.statistics.Settings.EnergyRange.Max()
			u.Balance = s.Config.InitialBalance
			s.model.InsertUser(u)
			return u, nil
		}
		return uu, err
	}
	return uu, nil
}

func (s *Service) FindUser(u model.User) (model.User, error) {
	u, err := s.model.GetUser(u)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return u, ErrUserNotFound
		}
		return u, err
	}
	return u, nil
}

func (s *Service) RaiseLimit(u model.User, l int) int {
	if u.DebtLimit < l {
		s.model.RaiseLimit(u, l)
		return l
	}
	return u.DebtLimit
}

func (s *Service) IncrementMessages(u model.User) {
	s.model.IncrementMessages(u)
}
