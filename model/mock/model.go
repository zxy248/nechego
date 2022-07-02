package mock

import "nechego/model"

func NewModel() *model.Model {
	return &model.Model{
		Admins:    &Admins{},
		Bans:      &Bans{},
		Eblans:    &Eblans{},
		Forbid:    &Forbid{},
		Pairs:     &Pairs{},
		Status:    &Status{},
		Users:     &Users{},
		Whitelist: &Whitelist{},
	}
}
