package model

import (
	"database/sql"
	"errors"
)

var ErrNoEblan = errors.New("no eblan")

type Eblans struct {
	DB *DB
}

const insertEblansQuery = `
insert into eblans (gid, uid, added)
values (?, ?, datetime('now', 'localtime'))
`

// Insert inserts the user to the eblans table.
func (e *Eblans) Insert(gid, uid int64) error {
	_, err := e.DB.Exec(insertEblansQuery, gid, uid)
	if err != nil {
		return err
	}
	return nil
}

const getEblansQuery = `
select uid from eblans
where gid = ? and added > date('now', 'localtime')
order by added desc
limit 1
`

// Get gets the user from the eblans table.
func (e *Eblans) Get(gid int64) (int64, error) {
	var uid int64
	if err := e.DB.QueryRow(getEblansQuery, gid).Scan(&uid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoEblan
		}
		return 0, err
	}
	return uid, nil
}

const deleteEblansQuery = `delete from eblans where gid = ?`

func (e *Eblans) Delete(gid int64) error {
	_, err := e.DB.Exec(deleteEblansQuery, gid)
	if err != nil {
		return err
	}
	return nil
}
