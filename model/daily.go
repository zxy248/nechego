package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Daily struct {
	ID     int
	GID    int64
	UserID int `db:"user_id"`
	Added  time.Time
}

const getDaily = `
select id, gid, user_id, added
from %s
where gid = ? and added > date('now', 'localtime')
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

func (m *Model) dailyUser(g Group, tableName string) (User, error) {
	var daily Daily
	var user User
	tx := m.db.MustBegin()
	defer tx.Rollback()

	if err := tx.Get(&daily, getDailyQuery(tableName), g.GID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// randomly choosing and inserting a new user if not found
			if err := tx.Get(&user, randomUser, g.GID); err != nil {
				return user, err
			}
			tx.MustExec(insertDailyQuery(tableName), user.GID, user.ID)
			return user, tx.Commit()
		}
		return user, err
	}
	if err := tx.Get(&user, getUserByID, daily.UserID); err != nil {
		return user, err
	}
	return user, tx.Commit()

}

func (m *Model) GetDailyEblan(g Group) (User, error) {
	return m.dailyUser(g, "daily_eblans")
}

func (m *Model) GetDailyAdmin(g Group) (User, error) {
	return m.dailyUser(g, "daily_admins")
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
where gid = ? and added > date('now', 'localtime')
limit 1`

const insertDailyPair = `
insert into daily_pairs (gid, user_id_x, user_id_y, added)
values (?, ?, ?, datetime('now', 'localtime'))`

func (m *Model) GetDailyPair(g Group) (User, User, error) {
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
