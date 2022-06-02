package main

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type pair struct {
	x int64
	y int64
}

var errNoPair = errors.New("no pair")
var errNoEblan = errors.New("no eblan")

type store struct {
	db *sql.DB
}

// newStore initializes the SQLite database and returns the store
func newStore(dsn string) (*store, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return &store{db}, nil
}

// insertUserID inserts the user ID to the users table
func (s *store) insertUserID(groupID, userID int64) error {
	query := "insert into users (user_id, group_id) values (?, ?)"
	_, err := s.db.Exec(query, userID, groupID)
	if err != nil {
		return err
	}
	return nil
}

// getUserIDs gets all the user IDs from the users table
func (s *store) getUserIDs(groupID int64) ([]int64, error) {
	query := "select user_id from users where group_id = ?"
	rows, err := s.db.Query(query, groupID)
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

// insertPair inserts the pair to the pairs table
func (s *store) insertPair(groupID int64, p pair) error {
	query := "insert into pairs (group_id, user_id_x, user_id_y, last) values (?, ?, ?, datetime('now', 'localtime'))"
	_, err := s.db.Exec(query, groupID, p.x, p.y)
	if err != nil {
		return err
	}
	return nil
}

// getPair gets the pair from the pairs table
func (s *store) getPair(groupID int64) (pair, error) {
	query := "select user_id_x, user_id_y from pairs where group_id = ? and last > date('now', 'localtime') order by last desc limit 1"
	var p pair
	if err := s.db.QueryRow(query, groupID).Scan(&p.x, &p.y); err != nil {
		if err == sql.ErrNoRows {
			return p, errNoPair
		}
		return p, err
	}
	return p, nil
}

// insertEblan inserts the user to the eblans table
func (s *store) insertEblan(groupID, userID int64) error {
	query := "insert into eblans (group_id, user_id, last) values (?, ?, datetime('now', 'localtime'))"
	_, err := s.db.Exec(query, groupID, userID)
	if err != nil {
		return err
	}
	return nil
}

// getEblan gets the user from the eblans table
func (s *store) getEblan(groupID int64) (int64, error) {
	query := "select user_id from eblans where group_id = ? and last > date('now', 'localtime') order by last desc limit 1"
	var userID int64
	if err := s.db.QueryRow(query, groupID).Scan(&userID); err != nil {
		if err == sql.ErrNoRows {
			return 0, errNoEblan
		}
		return 0, err
	}
	return userID, nil
}
