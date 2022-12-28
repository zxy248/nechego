package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type DailyType int

const (
	DailyEblan DailyType = iota
	DailyAdmin
)

var dailyTableNames = map[DailyType]string{
	DailyEblan: "daily_eblans",
	DailyAdmin: "daily_admins",
}

type Daily struct {
	ID     int
	GID    int64
	UserID int `db:"user_id"`
	Added  time.Time
}

const getDaily = `
select id, gid, user_id, added
from %s
where gid = ? and added >= date('now', 'localtime')
limit 1
`

func getDailyQuery(tableName string) string {
	return fmt.Sprintf(getDaily, tableName)
}

const insertDaily = `
insert into %s (gid, user_id, added)
values (?, ?, datetime('now', 'localtime'))`

func insertDailyQuery(tableName string) string {
	return fmt.Sprintf(insertDaily, tableName)
}

func (m *Model) dailyUser(g Group, d DailyType, roll bool, u User) (User, error) {
	tableName := dailyTableNames[d]
	var daily Daily
	tx := m.db.MustBegin()
	defer tx.Rollback()
	if err := tx.Get(&daily, getDailyQuery(tableName), g.GID); err != nil {
		if errors.Is(err, sql.ErrNoRows) && roll {
			if !u.Exists() {
				if err := tx.Get(&u, randomUser, g.GID); err != nil {
					return User{}, err
				}
			}
			tx.MustExec(insertDailyQuery(tableName), u.GID, u.ID)
			return u, tx.Commit()
		}
		return User{}, err
	}
	if err := tx.Get(&u, getUserByID, daily.UserID); err != nil {
		return User{}, err
	}
	return u, tx.Commit()
}

func (m *Model) DailyUser(g Group, d DailyType) (User, error) {
	return m.dailyUser(g, d, true, User{})
}

func (m *Model) DailyUserSet(g Group, d DailyType, u User) (User, error) {
	return m.dailyUser(g, d, true, u)
}

func (m *Model) DailyUserIfExists(g Group, d DailyType) (User, error) {
	u, err := m.dailyUser(g, d, false, User{})
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, nil
	}
	return u, err
}

type DailyPair struct {
	ID      int
	GID     int64
	UserIDX int `db:"user_id_x"`
	UserIDY int `db:"user_id_y"`
	Added   time.Time
}

const getDailyPair = `
select id, gid, user_id_x, user_id_y, added
from daily_pairs
where gid = ? and added >= date('now', 'localtime')
limit 1`

const insertDailyPair = `
insert into daily_pairs (gid, user_id_x, user_id_y, added)
values (?, ?, ?, datetime('now', 'localtime'))`

func (m *Model) DailyPair(g Group) (User, User, error) {
	var user1, user2 User
	var pair DailyPair
	tx := m.db.MustBegin()
	defer tx.Rollback()

	err := tx.Get(&pair, getDailyPair, g.GID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			var users []User
			if err := m.db.Select(&users, randomUsers, g.GID, 2); err != nil {
				return user1, user2, err
			}
			if len(users) != 2 {
				return user1, user2, errors.New("not enough users for a daily pair")
			}
			tx.MustExec(insertDailyPair, g.GID, users[0].ID, users[1].ID)
			return users[0], users[1], tx.Commit()
		}
		return user1, user2, err
	}
	if err := tx.Get(&user1, getUserByID, pair.UserIDX); err != nil {
		return user1, user2, err
	}
	if err := tx.Get(&user2, getUserByID, pair.UserIDY); err != nil {
		return user1, user2, err
	}
	return user1, user2, tx.Commit()
}
