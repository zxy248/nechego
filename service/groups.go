package service

import (
	"errors"
	"nechego/model"
)

func (s *Service) Group(g model.Group) (model.Group, error) {
	gg, err := s.model.GetGroup(g)
	if err != nil {
		if errors.Is(err, model.ErrGroupNotFound) {
			gg := model.Group{
				GID:         g.GID,
				Whitelisted: false,
				Status:      true,
			}
			s.model.InsertGroup(gg)
			return gg, nil
		}
		return gg, err
	}
	return gg, nil
}
