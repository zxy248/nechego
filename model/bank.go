package model

import "errors"

const deposit = `
update users set account = account + ?, balance = balance - ?
where id = ? and balance >= ?`

func (m *Model) Deposit(u User, amount, fee int) error {
	if amount <= 0 || fee < 0 {
		return ErrIncorrectAmount
	}
	n, err := m.db.MustExec(deposit, amount, amount+fee, u.ID, amount+fee).RowsAffected()
	failOn(err)
	if n != 1 {
		return ErrNotEnoughMoney
	}
	return nil
}

const withdraw = `
update users set account = account - ?, balance = balance + ?
where id = ? and account >= ?`

func (m *Model) Withdraw(u User, amount, fee int) error {
	if amount <= 0 || fee < 0 {
		return ErrIncorrectAmount
	}
	n, err := m.db.MustExec(withdraw, amount+fee, amount, u.ID, amount+fee).RowsAffected()
	failOn(err)
	if n != 1 {
		return ErrNotEnoughMoney
	}
	return nil
}

var ErrDebtLimit = errors.New("debt limit too low")

const debt = `
update users set balance = balance + ?, debt = debt + ?
where id = ? and debt = 0 and debt_limit >= ?`

func (m *Model) Debt(u User, amount, fee int) error {
	n, err := m.db.MustExec(debt, amount, amount+fee, u.ID, amount).RowsAffected()
	failOn(err)
	if n != 1 {
		return ErrDebtLimit
	}
	return nil
}

const repay = `
update users set account = account - ?, debt = debt - ?
where id = ? and account >= ? and debt >= ?`

func (m *Model) Repay(u User, amount int) error {
	n, err := m.db.MustExec(repay, amount, amount, u.ID, amount, amount).RowsAffected()
	failOn(err)
	if n != 1 {
		return ErrNotEnoughMoney
	}
	return nil
}

func (m *Model) DepositsToday(u User) (int, error) {
	var c int
	err := m.db.Get(&c, countUserEventsToday, depositEvent, u.GID, u.ID)
	return c, err
}

func (m *Model) AddDeposit(u User) {
	m.db.MustExec(insertEvent, u.GID, u.ID, depositEvent)
}
