package service

import (
	"errors"
	"nechego/input"
	"nechego/model"
	"nechego/statistics"
)

var (
	ErrAlreadyTurnedOn  = errors.New("already on")
	ErrAlreadyTurnedOff = errors.New("already off")
)

func (s *Service) TurnOn(g model.Group) error {
	if ok := s.model.EnableGroup(g); !ok {
		return ErrAlreadyTurnedOn
	}
	return nil
}

func (s *Service) TurnOff(g model.Group) error {
	if ok := s.model.DisableGroup(g); !ok {
		return ErrAlreadyTurnedOff
	}
	return nil
}

func (s *Service) Admins(g model.Group) ([]model.User, error) {
	return s.statistics.FilterUsers(g, statistics.IsAdmin)
}

func (s *Service) Bans(g model.Group) ([]model.User, error) {
	return s.statistics.FilterUsers(g, statistics.IsBanned)
}

func (s *Service) ForbiddenCommands(g model.Group) ([]input.Command, error) {
	return s.model.ForbiddenCommands(g)
}

func (s *Service) IsCommandForbidden(g model.Group, c input.Command) (bool, error) {
	return s.model.IsCommandForbidden(g, c)
}

func (s *Service) DeleteUsers(g model.Group, p func(model.User) bool) error {
	users, err := s.model.ListUsers(g)
	if err != nil {
		return err
	}
	for _, u := range users {
		if p(u) {
			s.model.DeleteUser(u)
		}
	}
	return nil
}
