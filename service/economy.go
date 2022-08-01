package service

import (
	"errors"
	"fmt"
	"nechego/model"
)

var (
	ErrNotEnoughMoney  = errors.New("not enough money")
	ErrIncorrectAmount = errors.New("incorrect amount")
)

type NotEnoughMoneyError struct {
	Delta int
}

func (e NotEnoughMoneyError) Error() string {
	return fmt.Sprintf("not enough money (%d)", e.Delta)
}

func (s *Service) Transfer(sender, recipient model.User, amount int) error {
	if err := s.model.TransferMoney(sender, recipient, amount); err != nil {
		if errors.Is(err, model.ErrNotEnoughMoney) {
			return NotEnoughMoneyError{amount - sender.Balance}
		}
		return err
	}
	return nil
}
