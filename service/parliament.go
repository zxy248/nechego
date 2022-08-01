package service

import (
	"errors"
	"nechego/model"
)

func (s *Service) Parliament(g model.Group) ([]model.User, error) {
	parliament, err := s.model.Parliament(g, s.Config.ParliamentMembers)
	if err != nil {
		return nil, err
	}
	return parliament, nil
}

var (
	ErrNotParliamentMember = errors.New("not a parliament member")
	ErrAlreadyVoted        = errors.New("already voted")
	ErrAlreadyImpeached    = errors.New("already impeached")
)

func (s *Service) Impeachment(g model.Group, u model.User) (votesLeft int, err error) {
	votes, err := s.model.Impeachment(g, u, s.Config.ParliamentMajority)
	if err != nil {
		if errors.Is(err, model.ErrNotParliamentMember) {
			return 0, ErrNotParliamentMember
		}
		if errors.Is(err, model.ErrAlreadyVoted) {
			return 0, ErrAlreadyVoted
		}
		if errors.Is(err, model.ErrAlreadyImpeached) {
			return 0, ErrAlreadyImpeached
		}
		return 0, err
	}
	return s.Config.ParliamentMajority - votes, nil
}
