package sqlite

import (
	"database/sql"
	"errors"
	"nechego/model"
)

type Economy struct {
	DB *DB
}

const (
	economyVerifyAmountQuery    = "select balance >= ? from users where gid = ? and uid = ?"
	economyVerifyRecipientQuery = "select 1 from users where gid = ? and uid = ?"
	economyUpdateBalanceQuery   = "update users set balance = balance + ? where gid = ? and uid = ?"
)

// Transfer sends the specified amount of money from one user to another.
func (e *Economy) Transfer(gid, sender, recipient int64, amount uint) error {
	tx, err := e.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var i int
	if err := tx.QueryRow(economyVerifyRecipientQuery, gid, recipient).Scan(&i); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrNoUser
		}
		return err
	}

	var enough bool
	if err := tx.QueryRow(economyVerifyAmountQuery, amount, gid, sender).Scan(&enough); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrNoUser
		}
		return err
	}
	if !enough {
		return model.ErrNotEnoughMoney
	}

	_, err = tx.Exec(economyUpdateBalanceQuery, -amount, gid, sender)
	if err != nil {
		return err
	}
	_, err = tx.Exec(economyUpdateBalanceQuery, amount, gid, recipient)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

const economyBalanceQuery = "select balance from users where gid = ? and uid = ?"

// Balance returns an amount of money on the user's balance.
func (e *Economy) Balance(gid, uid int64) (uint, error) {
	var amount uint
	if err := e.DB.QueryRow(economyBalanceQuery, gid, uid).Scan(&amount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, model.ErrNoUser
		}
		return 0, err
	}
	return amount, nil
}

// Update changes the user's balance by delta.
// Returns model.ErrNotEnoughMoney and does nothing when the user's balance would be less than 0.
func (e *Economy) Update(gid, uid int64, delta int) error {
	tx, err := e.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if delta < 0 {
		var enough bool
		if err := tx.QueryRow(economyVerifyAmountQuery, -delta, gid, uid).Scan(&enough); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return model.ErrNoUser
			}
			return err
		}
		if !enough {
			return model.ErrNotEnoughMoney
		}
	}

	_, err = tx.Exec(economyUpdateBalanceQuery, delta, gid, uid)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
