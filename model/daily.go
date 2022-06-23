package model

import (
	"database/sql"
	"fmt"
)

const insertDailyQuery = `
insert into %s (gid, uid, added)
values (?, ?, datetime('now', 'localtime'))
`

func insertDaily(db *sql.DB, table string, gid, uid int64) error {
	_, err := db.Exec(fmt.Sprintf(insertDailyQuery, table), gid, uid)
	return err
}

const getDailyQuery = `
select uid from %s
where gid = ? and added > date('now', 'localtime')
order by added desc
limit 1
`

func getDaily(db *sql.DB, table string, gid int64) (int64, error) {
	var uid int64
	if err := db.QueryRow(fmt.Sprintf(getDailyQuery, table), gid).Scan(&uid); err != nil {
		return 0, err
	}
	return uid, nil
}

const deleteDailyQuery = `delete from %s where gid = ?`

func deleteDaily(db *sql.DB, table string, gid int64) error {
	_, err := db.Exec(fmt.Sprintf(deleteDailyQuery, table), gid)
	return err
}
