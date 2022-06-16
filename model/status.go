package model

import (
	"database/sql"
	"errors"
)

type Status struct {
	DB *DB
}

const enableStatusQuery = `delete from status where gid = ?`

func (s *Status) Enable(gid int64) error {
	_, err := s.DB.Exec(enableStatusQuery, gid)
	if err != nil {
		return err
	}
	return nil
}

const activeStatusQuery = `select 1 from status where gid = ?`

func (s *Status) Active(gid int64) (bool, error) {
	var i int
	if err := s.DB.QueryRow(activeStatusQuery, gid).Scan(&i); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, err
	}
	return !(i == 1), nil
}

const disableStatusQuery = `insert into status (gid) values (?)`

func (s *Status) Disable(gid int64) error {
	_, err := s.DB.Exec(disableStatusQuery, gid)
	if err != nil {
		return err
	}
	return nil
}
