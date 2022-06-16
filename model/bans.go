package model

import (
	"database/sql"
	"errors"
)

type Bans struct {
	DB *DB
}

const insertBansQuery = `insert into bans (uid) values (?)`

func (b *Bans) Ban(uid int64) error {
	_, err := b.DB.Exec(insertBansQuery, uid)
	if err != nil {
		return err
	}
	return nil
}

const deleteBansQuery = `delete from bans where uid = ?`

func (b *Bans) Unban(uid int64) error {
	_, err := b.DB.Exec(deleteBansQuery, uid)
	if err != nil {
		return err
	}
	return nil
}

const listBansQuery = `select uid from bans`

// List returns a list of banned users.
func (b *Bans) List() ([]int64, error) {
	rows, err := b.DB.Query(listBansQuery)
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

const bannedBansQuery = `select 1 from bans where uid = ?`

func (b *Bans) Banned(uid int64) (bool, error) {
	var i int
	if err := b.DB.QueryRow(bannedBansQuery, uid).Scan(&i); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return i == 1, nil
}
