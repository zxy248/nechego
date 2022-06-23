package model

import (
	"database/sql"
	"errors"
)

var ErrNoEblan = errors.New("no eblan")

type Eblans struct {
	DB *DB
}

const eblansTableName = "eblans"

// Insert adds a user to the eblans table.
func (e *Eblans) Insert(gid, uid int64) error {
	return insertDaily(e.DB.DB, eblansTableName, gid, uid)
}

// Get returns the user ID of the eblan of the day. If there is no eblan, returns 0, ErrNoEblan.
func (e *Eblans) Get(gid int64) (int64, error) {
	uid, err := getDaily(e.DB.DB, eblansTableName, gid)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNoEblan
	}
	return uid, err
}

// Delete removes the current eblan of the day and all the previous.
func (e *Eblans) Delete(gid int64) error {
	return deleteDaily(e.DB.DB, eblansTableName, gid)
}
