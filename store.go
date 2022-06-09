package main

import (
	"database/sql"
	"errors"
)

type pairOfTheDay struct {
	userIDx int64
	userIDy int64
}

var errNoPair = errors.New("no pair")
var errNoEblan = errors.New("no eblan")

type store struct {
	db *sql.DB
}

// newStore initializes the SQLite database and returns the new store.
func newStore(dsn string) (*store, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return &store{db}, nil
}

const insertUserIDQuery = `
insert into users (user_id, group_id)
values (?, ?)
`

// insertUserID inserts the user ID to the users table.
func (s *store) insertUserID(groupID, userID int64) error {
	_, err := s.db.Exec(insertUserIDQuery, userID, groupID)
	if err != nil {
		return err
	}
	return nil
}

const getUserIDsQuery = `
select user_id from users
where group_id = ?
`

// getUserIDs gets all the user IDs from the users table.
func (s *store) getUserIDs(groupID int64) ([]int64, error) {
	rows, err := s.db.Query(getUserIDsQuery, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	userIDs := []int64{}
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		userIDs = append(userIDs, id)
	}
	return userIDs, nil
}

const insertPairQuery = `
insert into pairs (group_id, user_id_x, user_id_y, last)
values (?, ?, ?, datetime('now', 'localtime'))
`

// insertPair inserts the pair to the pairs table.
func (s *store) insertPair(groupID int64, p pairOfTheDay) error {
	_, err := s.db.Exec(insertPairQuery, groupID, p.userIDx, p.userIDy)
	if err != nil {
		return err
	}
	return nil
}

const getPairQuery = `
select user_id_x, user_id_y from pairs
where group_id = ? and last > date('now', 'localtime')
order by last desc
limit 1
`

// getPair gets the pair from the pairs table.
func (s *store) getPair(groupID int64) (pairOfTheDay, error) {
	var p pairOfTheDay
	if err := s.db.QueryRow(getPairQuery, groupID).Scan(&p.userIDx, &p.userIDy); err != nil {
		if err == sql.ErrNoRows {
			return p, errNoPair
		}
		return p, err
	}
	return p, nil
}

const insertEblanQuery = `
insert into eblans (group_id, user_id, last)
values (?, ?, datetime('now', 'localtime'))
`

// insertEblan inserts the user to the eblans table.
func (s *store) insertEblan(groupID, userID int64) error {
	_, err := s.db.Exec(insertEblanQuery, groupID, userID)
	if err != nil {
		return err
	}
	return nil
}

const getEblanQuery = `
select user_id from eblans
where group_id = ? and last > date('now', 'localtime')
order by last desc
limit 1
`

// getEblan gets the user from the eblans table.
func (s *store) getEblan(groupID int64) (int64, error) {
	var userID int64
	if err := s.db.QueryRow(getEblanQuery, groupID).Scan(&userID); err != nil {
		if err == sql.ErrNoRows {
			return 0, errNoEblan
		}
		return 0, err
	}
	return userID, nil
}
