package statistics

import (
	"errors"
	"nechego/model"
	"sort"
)

func (s *Statistics) SortedUsers(g model.Group, f UserSortFunc) ([]model.User, error) {
	users, err := s.model.ListUsers(g)
	if err != nil {
		return nil, err
	}
	sort.Slice(users, f(users))
	return users, nil
}

func (s *Statistics) GreatestUser(g model.Group, f UserSortFunc) (model.User, error) {
	users, err := s.SortedUsers(g, f)
	if err != nil {
		return model.User{}, err
	}
	if len(users) < 1 {
		return model.User{}, errors.New("empty list")
	}
	return users[0], nil
}

type UserSortFunc func([]model.User) func(int, int) bool

func ByWealthDesc(users []model.User) func(int, int) bool {
	return func(i, j int) bool {
		return users[i].Summary() > users[j].Summary()
	}
}

func ByWealthAsc(users []model.User) func(int, int) bool {
	return func(i, j int) bool {
		return users[i].Summary() < users[j].Summary()
	}
}

func ByEloDesc(users []model.User) func(int, int) bool {
	return func(i, j int) bool {
		return users[i].Elo > users[j].Elo
	}
}

func ByEloAsc(users []model.User) func(int, int) bool {
	return func(i, j int) bool {
		return users[i].Elo < users[j].Elo
	}
}

func (s *Statistics) ByStrengthDesc(users []model.User) func(int, int) bool {
	return func(i, j int) bool {
		ii, _ := s.Strength(users[i])
		jj, _ := s.Strength(users[j])
		return ii > jj
	}
}

func (s *Statistics) ByStrengthAsc(users []model.User) func(int, int) bool {
	return func(i, j int) bool {
		ii, _ := s.Strength(users[i])
		jj, _ := s.Strength(users[j])
		return ii < jj
	}
}

func (s *Statistics) FilterUsers(g model.Group, f UserFilterFunc) ([]model.User, error) {
	users, err := s.model.ListUsers(g)
	if err != nil {
		return nil, err
	}
	var filtered []model.User
	for _, u := range users {
		if f(u) {
			filtered = append(filtered, u)
		}
	}
	return filtered, nil
}

type UserFilterFunc func(model.User) bool

func IsAdmin(u model.User) bool {
	return u.Admin
}

func IsBanned(u model.User) bool {
	return u.Banned
}

func (s *Statistics) IsPoor(u model.User) bool {
	return u.Summary() < s.Settings.PoorThreshold
}

func (s *Statistics) IsRich(u model.User) (bool, error) {
	rich, err := s.GreatestUser(model.Group{GID: u.GID}, ByWealthDesc)
	if err != nil {
		return false, err
	}
	return rich.ID == u.ID, nil
}
