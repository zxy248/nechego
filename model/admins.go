package model

import (
	"database/sql"
	"errors"
)

type Admins struct {
	DB *DB
}

const insertAdminsQuery = `insert into admins (uid) values (?)`

// Insert inserts a user ID to the admins table.
func (a *Admins) Insert(uid int64) error {
	_, err := a.DB.Exec(insertAdminsQuery, uid)
	if err != nil {
		return err
	}
	return nil
}

const deleteAdminsQuery = `delete from admins where uid = ?`

func (a *Admins) Delete(uid int64) error {
	_, err := a.DB.Exec(deleteAdminsQuery, uid)
	if err != nil {
		return err
	}
	return nil
}

const listAdminsQuery = `select uid from admins`

// List returns the list of admins.
func (a *Admins) List() ([]int64, error) {
	rows, err := a.DB.Query(listAdminsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

const authorizeAdminsQuery = `select 1 from admins where uid = ?`

func (a *Admins) Authorize(uid int64) (bool, error) {
	var i int
	if err := a.DB.QueryRow(authorizeAdminsQuery, uid).Scan(&i); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return i == 1, nil
}
