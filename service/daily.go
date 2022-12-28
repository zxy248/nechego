package service

import (
	"nechego/model"
	"time"
)

func (s *Service) DailyPair(g model.Group, u model.User) (model.User, model.User, error) {
	return s.model.DailyPair(g)
}

func (s *Service) DailyEblan(g model.Group, u model.User) (model.User, error) {
	if time.Now().Weekday() == time.Friday {
		return s.model.DailyUserSet(g, model.DailyEblan, u)
	}
	return s.model.DailyUser(g, model.DailyEblan)
}

func (s *Service) DailyAdmin(g model.Group, u model.User) (model.User, error) {
	return s.model.DailyUser(g, model.DailyAdmin)
}
