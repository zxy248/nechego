package model

import (
	"database/sql"
	"errors"
)

type Whitelist struct {
	DB *DB
}

const insertWhitelistQuery = `insert into whitelist (gid) values (?)`

// Insert inserts a group ID to the whitelist table.
func (w *Whitelist) Insert(gid int64) error {
	_, err := w.DB.Exec(insertWhitelistQuery, gid)
	if err != nil {
		return err
	}
	return nil
}

const allowWhitelistQuery = `select 1 from whitelist where gid = ?`

// Allow returns true if the given group ID exists in the whitelist table.
func (w *Whitelist) Allow(gid int64) (bool, error) {
	var i int
	if err := w.DB.QueryRow(allowWhitelistQuery, gid).Scan(&i); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return i == 1, nil
}
