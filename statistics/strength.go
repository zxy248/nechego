package statistics

import "nechego/model"

func (s *Statistics) Strength(u model.User) (float64, error) {
	mcc, err := s.messageCountCoefficient(u)
	if err != nil {
		return 0, err
	}
	mul, err := s.strengthMultiplier(u)
	if err != nil {
		return 0, err
	}
	return (1.0 + mcc) * mul, nil
}

func (s *Statistics) messageCountCoefficient(u model.User) (float64, error) {
	c, err := s.model.GroupMessageCount(model.Group{GID: u.GID})
	if err != nil {
		return 0, err
	}
	return (1.0 + float64(u.Messages)) / (1.0 + float64(c)), nil
}

func (s *Statistics) strengthMultiplier(u model.User) (float64, error) {
	m, err := s.UserModset(u)
	if err != nil {
		return 0, err
	}
	return 1.0 + m.Sum(), nil
}
