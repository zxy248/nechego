package model

import (
	"database/sql"
	"errors"
)

type User struct {
	ID       int
	GID      int64
	UID      int64
	Energy   int
	Balance  int
	Admin    bool
	Banned   bool
	Messages int
}

const insertUser = `
insert into users (gid, uid, energy, balance, admin, banned, messages)
values (?, ?, ?, ?, ?, ?, ?)`

func (m *Model) InsertUser(u User) {
	m.db.MustExec(insertUser,
		u.GID, u.UID, u.Energy, u.Balance, u.Admin, u.Banned, u.Messages)
}

const deleteUser = `
delete from users
where gid = ? and uid = ?`

func (m *Model) DeleteUser(u User) {
	m.db.MustExec(deleteUser, u.GID, u.UID)
}

const selectUser = `
select id, gid, uid, energy, balance, admin, banned, messages
from users`

const (
	userByID   = "id = ?"
	userByGUID = "(gid = ? and uid = ?)"
)

var (
	getUserByID   = concat(selectUser, "where", userByID, "limit 1")
	getUserByGUID = concat(selectUser, "where", userByGUID, "limit 1")
)

var ErrUserNotFound = errors.New("user not found")

func (m *Model) GetUser(u User) (User, error) {
	user := User{}
	err := m.db.Get(&user, getUserByGUID, u.GID, u.UID)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrUserNotFound
	}
	return user, err
}

const listUsers = selectUser + `
where gid = ?`

func (m *Model) ListUsers(g Group) ([]User, error) {
	users := []User{}
	err := m.db.Select(&users, listUsers, g.GID)
	return users, err
}

const randomUser = selectUser + `
where gid = ?
order by random()
limit 1`

func (m *Model) RandomUser(g Group) (User, error) {
	user := User{}
	err := m.db.Get(&user, randomUser, g.GID)
	return user, err
}

const randomUsers = selectUser + `
where gid = ?
order by random()
limit ?`

func (m *Model) RandomUsers(g Group, n int) ([]User, error) {
	users := []User{}
	err := m.db.Select(&users, randomUsers, g.GID, n)
	return users, err
}

const banUser = `
update users set banned = ?
where gid = ? and uid = ?`

func (m *Model) BanUser(u User) {
	m.db.MustExec(banUser, true, u.GID, u.UID)
}

func (m *Model) UnbanUser(u User) {
	m.db.MustExec(banUser, false, u.GID, u.UID)
}

const restoreEnergy = `
update users set energy = energy + ?
where energy + ? <= ?`

func (m *Model) RestoreEnergy(delta, cap int) {
	m.db.MustExec(restoreEnergy, delta, delta, cap)
}

const updateEnergy = `
update users set energy = energy + ?
where id = ? and (energy + ? <= ?) and (energy + ? >= 0)`

func (m *Model) UpdateEnergy(u User, delta, cap int) (updated bool) {
	n, err := m.db.MustExec(updateEnergy, delta, u.ID, delta, cap, delta).RowsAffected()
	failOn(err)
	return n == 1
}

const updateBalance = `
update users set balance = balance + ?
where id = ? and (balance + ? >= 0)`

var (
	ErrNotEnoughMoney = errors.New("sender has not enough money")
)

func (m *Model) TransferMoney(sender, recipient User, amount int) error {
	tx := m.db.MustBegin()
	defer tx.Rollback()
	n, err := tx.MustExec(updateBalance, -amount, sender.ID, -amount).RowsAffected()
	failOn(err)
	if n != 1 {
		return ErrNotEnoughMoney
	}
	n, err = tx.MustExec(updateBalance, +amount, recipient.ID, +amount).RowsAffected()
	failOn(err)
	if n != 1 {
		return ErrNotEnoughMoney
	}
	return tx.Commit()
}

func (m *Model) UpdateMoney(u User, amount int) (updated bool) {
	n, err := m.db.MustExec(updateBalance, +amount, u.ID, +amount).RowsAffected()
	failOn(err)
	return n == 1
}

const incrementMessages = `
update users set messages = messages + 1
where id = ?`

func (m *Model) IncrementMessages(u User) {
	m.db.MustExec(incrementMessages, u.ID)
}
