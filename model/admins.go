package model

import (
	"database/sql"
	"errors"

	"golang.org/x/exp/slices"
)

type Admins struct {
	DB *DB
}

const insertAdminsQuery = `insert into admins (uid) values (?)`

// Insert adds an admin.
func (a *Admins) Insert(uid int64) error {
	return insertAdmin(a.DB.DB, uid)
}

const deleteAdminsQuery = `delete from admins where uid = ?`

// Delete removes an admin.
func (a *Admins) Delete(uid int64) error {
	return deleteAdmin(a.DB.DB, uid)
}

const listAdminsQuery = `select uid from admins`

// List returns a list of admins.
func (a *Admins) List(gid int64) ([]int64, error) {
	l, err := listAdmins(a.DB.DB)
	if err != nil {
		return nil, err
	}
	i, err := a.GetDaily(gid)
	if err != nil {
		return l, nil
	}
	if slices.Contains(l, i) {
		return []int64{i}, nil
	}
	return append(l, i), nil
}

const authorizeAdminsQuery = `select 1 from admins where uid = ?`

// Authorize returns true if the given user is an admin, false otherwise.
func (a *Admins) Authorize(gid, uid int64) (bool, error) {
	l, err := a.List(gid)
	if err != nil {
		return false, err
	}
	return slices.Contains(l, uid), nil
}

func insertAdmin(db *sql.DB, uid int64) error {
	_, err := db.Exec(insertAdminsQuery, uid)
	return err
}

func deleteAdmin(db *sql.DB, uid int64) error {
	_, err := db.Exec(deleteAdminsQuery, uid)
	return err
}

func listAdmins(db *sql.DB) ([]int64, error) {
	rows, err := db.Query(listAdminsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
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

func isAdmin(db *sql.DB, uid int64) (bool, error) {
	var i int
	if err := db.QueryRow(authorizeAdminsQuery, uid).Scan(&i); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

var ErrNoAdmin = errors.New("no admin")

const adminTableName = "admin"

func (a *Admins) InsertDaily(gid, uid int64) error {
	return insertDaily(a.DB.DB, adminTableName, gid, uid)
}

func (a *Admins) GetDaily(gid int64) (int64, error) {
	uid, err := getDaily(a.DB.DB, adminTableName, gid)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNoAdmin
	}
	return uid, err
}

func (a *Admins) DeleteDaily(gid int64) error {
	return deleteDaily(a.DB.DB, adminTableName, gid)
}
