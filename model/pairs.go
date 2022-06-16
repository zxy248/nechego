package model

import (
	"database/sql"
	"errors"
)

var ErrNoPair = errors.New("no pair")

type Pairs struct {
	DB *DB
}

const insertPairQuery = `
insert into pairs (gid, uidx, uidy, added)
values (?, ?, ?, datetime('now', 'localtime'))
`

// Insert inserts the pair to the pairs table.
func (p *Pairs) Insert(gid int64, uidx, uidy int64) error {
	_, err := p.DB.Exec(insertPairQuery, gid, uidx, uidy)
	if err != nil {
		return err
	}
	return nil
}

const getPairQuery = `
select uidx, uidy from pairs
where gid = ? and added > date('now', 'localtime')
order by added desc
limit 1
`

// Get gets the pair from the pairs table.
func (p *Pairs) Get(gid int64) (uidx, uidy int64, err error) {
	if err := p.DB.QueryRow(getPairQuery, gid).Scan(&uidx, &uidy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, ErrNoPair
		}
		return 0, 0, err
	}
	return uidx, uidy, nil
}
