package service

import (
	"errors"
	"nechego/fishing"
	"nechego/model"
	"nechego/numbers"
	"nechego/statistics"
)

var (
	ErrNotHungry     = errors.New("not hungry")
	ErrNotEnoughFish = errors.New("not enough fish")
	ErrNotFisher     = errors.New("not a fisher")
	ErrAlreadyFisher = errors.New("already a fisher")
)

func (s *Service) EatFish(u model.User) (energyRestored int, err error) {
	delta := numbers.InRange(
		s.Config.EatEnergyRestore.Min(),
		s.Config.EatEnergyRestore.Max(),
	)
	if s.statistics.HasFullEnergy(u) {
		return 0, ErrNotHungry
	}
	enoughFish := s.model.EatFish(u, delta, statistics.EnergyHardLimit)
	if !enoughFish {
		return 0, ErrNotEnoughFish
	}
	return delta, nil
}

func (s *Service) BuyFishingRod(u model.User) error {
	if u.Fisher {
		return ErrAlreadyFisher
	}
	enough := s.model.UpdateMoney(u, -s.Config.FishingRodPrice)
	if !enough {
		return NotEnoughMoneyError{s.Config.FishingRodPrice - u.Balance}
	}
	s.model.AllowFishing(u)
	return nil
}

func (s *Service) Fish(u model.User) (fishing.Session, error) {
	if !u.Fisher {
		return fishing.Session{}, ErrNotFisher
	}
	enoughEnergy := s.model.UpdateEnergy(u, -s.Config.FishingEnergyDrain, statistics.EnergyHardLimit)
	if !enoughEnergy {
		return fishing.Session{}, ErrNotEnoughEnergy
	}
	session := fishing.Cast()
	if session.Success() {
		s.collectFish(u, session.Fish)
	}
	return session, nil
}

func (s *Service) collectFish(u model.User, f fishing.Fish) {
	if f.Light() {
		s.model.AddFish(u)
		return
	}
	s.model.InsertFish(model.MakeCatch(u, f))
}

func (s *Service) FreshFish(u model.User) (fishing.Fishes, error) {
	catch, err := s.model.SelectFish(u)
	if err != nil {
		return nil, err
	}
	fishes := fishing.Fishes{}
	for _, c := range catch {
		if c.Frozen {
			continue
		}
		fishes = append(fishes, c.Fish)
	}
	return fishes, nil
}

func (s *Service) SellFish(u model.User) (fishing.Fishes, error) {
	catch, err := s.model.SellFish(u)
	if err != nil {
		return nil, err
	}
	fishes := fishing.Fishes{}
	for _, c := range catch {
		fishes = append(fishes, c.Fish)
	}
	s.model.UpdateMoney(u, fishes.Price())
	return fishes, nil
}

func (s *Service) FreezeFish(u model.User) {
	s.model.FreezeFish(u)
}

func (s *Service) UnfreezeFish(u model.User) {
	s.model.UnfreezeFish(u)
}

func (s *Service) Freezer(u model.User) (fishing.Fishes, error) {
	catch, err := s.model.SelectFish(u)
	if err != nil {
		return nil, err
	}
	fishes := fishing.Fishes{}
	for _, c := range catch {
		if !c.Frozen {
			continue
		}
		fishes = append(fishes, c.Fish)
	}
	return fishes, nil

}
