package service

import (
	"errors"
	"math/rand"
	"nechego/model"
	"nechego/pets"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrPetAlreadyNamed = errors.New("pet already named")
	ErrPetAlreadyTaken = errors.New("pet already taken")
	ErrPetBadName      = errors.New("bad pet name")
	ErrNoPet           = errors.New("no pet")
)

func (s *Service) GetPet(u model.User) (*pets.Pet, error) {
	p, err := s.model.GetPet(u)
	if err != nil {
		if errors.Is(err, model.ErrNoPet) {
			return nil, ErrNoPet
		}
		return nil, err
	}
	return p, nil
}

func (s *Service) BuyPet(u model.User) (*pets.Pet, error) {
	hasPet, err := s.model.HasPet(u)
	if err != nil {
		return nil, err
	}
	if hasPet {
		return nil, ErrPetAlreadyTaken
	}
	if enough := s.model.UpdateMoney(u, -s.Config.PetPrice); !enough {
		return nil, NotEnoughMoneyError{s.Config.PetPrice - u.Balance}
	}
	p := pets.RandomPet(rand.Float64())
	s.model.InsertPet(u, p)
	return p, nil
}

func (s *Service) NamePet(u model.User, name string) error {
	if !validPetName(name) {
		return ErrPetBadName
	}
	p, err := s.model.GetPet(u)
	if err != nil {
		if errors.Is(err, model.ErrNoPet) {
			return ErrNoPet
		}
		return err
	}
	if p.HasName() {
		return ErrPetAlreadyNamed
	}
	s.model.NamePet(u, normalizePetName(name))
	return nil
}

func validPetName(s string) (ok bool) {
	return utf8.ValidString(s) && validPetNameLength(s) && validPetNameAlphabet(s)
}

func validPetNameAlphabet(s string) bool {
	for _, r := range s {
		if !(unicode.Is(unicode.Cyrillic, r) || r == ' ') {
			return false
		}
	}
	return true
}

func validPetNameLength(s string) bool {
	return utf8.RuneCountInString(s) < 32
}

func normalizePetName(s string) string {
	return strings.Title(s)
}

func (s *Service) DropPet(u model.User) error {
	hasPet, err := s.model.HasPet(u)
	if err != nil {
		return err
	}
	if !hasPet {
		return ErrNoPet
	}
	s.model.DeletePet(u)
	return nil
}
