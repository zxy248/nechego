package service

import (
	"errors"
	"nechego/model"
)

var (
	ErrMinDebt            = errors.New("debt too low")
	ErrDebtLimit          = errors.New("debt limit too low")
	ErrBankOperationLimit = errors.New("bank operation limit")
)

func (s *Service) Deposit(u model.User, amount int) (transfered int, err error) {
	c, err := s.model.DepositsToday(u)
	if err != nil {
		return 0, err
	}
	if c >= s.Config.MaxDeposits {
		return 0, ErrBankOperationLimit
	}
	taxed := amount - s.Config.DepositFee
	if err := s.model.Deposit(u, taxed, s.Config.DepositFee); err != nil {
		if errors.Is(err, model.ErrNotEnoughMoney) {
			return 0, NotEnoughMoneyError{amount - u.Balance}
		}
		if errors.Is(err, model.ErrIncorrectAmount) {
			return 0, ErrIncorrectAmount
		}
		return 0, err
	}
	s.model.AddDeposit(u)
	return taxed, nil
}

func (s *Service) Withdraw(u model.User, amount int) (transfered int, err error) {
	taxed := amount - s.Config.WithdrawFee
	if err := s.model.Withdraw(u, taxed, s.Config.WithdrawFee); err != nil {
		if errors.Is(err, model.ErrNotEnoughMoney) {
			return 0, NotEnoughMoneyError{amount - u.Account}
		}
		if errors.Is(err, model.ErrIncorrectAmount) {
			return 0, ErrIncorrectAmount
		}
		return 0, err
	}
	return taxed, nil
}

func (s *Service) Debt(u model.User, amount int) (debt int, err error) {
	if amount < s.Config.MinDebt {
		return 0, ErrMinDebt
	}
	fee := int(float64(amount)*s.Config.DebtPercentage) + 1
	if err := s.model.Debt(u, amount, fee); err != nil {
		if errors.Is(err, model.ErrDebtLimit) {
			return 0, ErrDebtLimit
		}
		return 0, err
	}
	return amount + fee, nil
}

func (s *Service) Repay(u model.User, amount int) (debt int, err error) {
	if u.Debt < amount {
		amount = u.Debt
	}
	if err := s.model.Repay(u, amount); err != nil {
		if errors.Is(err, model.ErrNotEnoughMoney) {
			return 0, NotEnoughMoneyError{amount - u.Account}
		}
		return 0, err
	}
	return u.Debt - amount, nil
}
