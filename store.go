package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

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
func (s *store) insertUserID(groupID int64, userID int64) error {
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
