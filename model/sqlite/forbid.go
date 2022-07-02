package sqlite

import (
	"database/sql"
	"errors"
	"nechego/input"
)

type Forbid struct {
	DB *DB
}

const forbidForbidQuery = `insert into forbid (gid, command) values (?, ?)`

func (f *Forbid) Forbid(gid int64, c input.Command) error {
	_, err := f.DB.Exec(forbidForbidQuery, gid, c)
	if err != nil {
		return err
	}
	return nil
}

const permitForbidQuery = `delete from forbid where gid = ? and command = ?`

func (f *Forbid) Permit(gid int64, c input.Command) error {
	_, err := f.DB.Exec(permitForbidQuery, gid, c)
	if err != nil {
		return err
	}
	return nil
}

const checkForbidQuery = `select 1 from forbid where gid = ? and command = ?`

func (f *Forbid) Forbidden(gid int64, c input.Command) (bool, error) {
	var i int
	if err := f.DB.QueryRow(checkForbidQuery, gid, c).Scan(&i); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

const listForbidQuery = `select command from forbid where gid = ?`

func (f *Forbid) List(gid int64) ([]input.Command, error) {
	rows, err := f.DB.Query(listForbidQuery, gid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var commands []input.Command
	for rows.Next() {
		var command input.Command
		if err := rows.Scan(&command); err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return commands, nil
}
