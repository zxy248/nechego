package service

import (
	"errors"
	"nechego/input"
	"nechego/model"
)

var (
	ErrAlreadyBanned    = errors.New("already banned")
	ErrAlreadyUnbanned  = errors.New("already unbanned")
	ErrAlreadyForbidden = errors.New("command already forbidden")
	ErrAlreadyPermitted = errors.New("command already permitted")
)

func (s *Service) Ban(u model.User) error {
	if u.Banned {
		return ErrAlreadyBanned
	}
	s.model.BanUser(u)
	return nil
}

func (s *Service) Unban(u model.User) error {
	if !u.Banned {
		return ErrAlreadyUnbanned
	}
	s.model.UnbanUser(u)
	return nil
}

func (s *Service) Forbid(g model.Group, c input.Command) error {
	if ok := s.model.ForbidCommand(g, c); !ok {
		return ErrAlreadyForbidden
	}
	return nil
}

func (s *Service) Permit(g model.Group, c input.Command) error {
	if ok := s.model.PermitCommand(g, c); !ok {
		return ErrAlreadyPermitted
	}
	return nil
}

func (s *Service) PermitAll(g model.Group) error {
	cmds, err := s.ForbiddenCommands(g)
	if err != nil {
		return err
	}
	for _, cmd := range cmds {
		if err := s.Permit(g, cmd); err != nil {
			return err
		}
	}
	return nil
}
