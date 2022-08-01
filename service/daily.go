package service

import "nechego/model"

func (s *Service) DailyPair(g model.Group) (model.User, model.User, error) {
	return s.model.DailyPair(g)
}

func (s *Service) DailyEblan(g model.Group) (model.User, error) {
	return s.model.DailyEblan(g)
}

func (s *Service) DailyAdmin(g model.Group) (model.User, error) {
	return s.model.DailyAdmin(g)
}
