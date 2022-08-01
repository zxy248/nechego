package statistics

import (
	"errors"
	"nechego/model"
)

func (s *Statistics) GroupBalance(g model.Group) (int, error) {
	users, err := s.model.ListUsers(g)
	if err != nil {
		return 0, err
	}
	return sumNetWorth(users), nil
}

func (s *Statistics) AverageBalance(g model.Group) (float64, error) {
	total, err := s.GroupBalance(g)
	if err != nil {
		return 0, err
	}
	c, err := s.model.CountUsers(g)
	if err != nil {
		return 0, err
	}
	if c == 0 {
		return 0, errors.New("no users")
	}
	return float64(total) / float64(c), nil
}

func sumNetWorth(u []model.User) int {
	sum := 0
	for _, uu := range u {
		sum += uu.Summary()
	}
	return sum
}
