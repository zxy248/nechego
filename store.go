package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type pair struct {
	x int64
	y int64
}

var errNoPair = errors.New("no pair")

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

// insertUserID inserts userID with groupID
func (s *store) insertUserID(groupID, userID int64) error {
	query := "insert into users (user_id, group_id) values (?, ?)"
	_, err := s.db.Exec(query, userID, groupID)
	if err != nil {
		return err
	}
	return nil
}

// getUserIDs gets all userIDs with groupID
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

func (s *store) insertPair(groupID int64, p pair) error {
	query := "insert into pairs (group_id, user_id_x, user_id_y, last) values (?, ?, ?, datetime())"
	_, err := s.db.Exec(query, groupID, p.x, p.y)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) getPair(groupID int64) (pair, error) {
	query := "select user_id_x, user_id_y from pairs where group_id = ? and last > date() order by last desc limit 1"
	var p pair
	if err := s.db.QueryRow(query, groupID).Scan(&p.x, &p.y); err != nil {
		fmt.Println(p)
		if err == sql.ErrNoRows {
			return p, errNoPair
		}
		return p, err
	}
	return p, nil
}
