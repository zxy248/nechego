package model

import (
	"database/sql"
	"errors"
)

type Group struct {
	GID         int64
	Whitelisted bool
	Status      bool
}

const insertGroup = `
insert into groups (gid, whitelisted, status)
values (?, ?, ?)`

func (m *Model) InsertGroup(g Group) {
	m.db.MustExec(insertGroup, g.GID, g.Whitelisted, g.Status)
}

const getGroup = `
select gid, whitelisted, status
from groups
where gid = ?`

var ErrGroupNotFound = errors.New("group not found")

func (m *Model) GetGroup(g Group) (Group, error) {
	group := Group{}
	err := m.db.Get(&group, getGroup, g.GID)
	if errors.Is(err, sql.ErrNoRows) {
		return group, ErrGroupNotFound
	}
	return group, err
}

const enableGroup = `
update groups set status = 1
where gid = ? and status = 0`

func (m *Model) EnableGroup(g Group) (updated bool) {
	n, err := m.db.MustExec(enableGroup, g.GID).RowsAffected()
	failOn(err)
	return n == 1
}

const disableGroup = `
update groups set status = 0
where gid = ? and status = 1`

func (m *Model) DisableGroup(g Group) (updated bool) {
	n, err := m.db.MustExec(disableGroup, g.GID).RowsAffected()
	failOn(err)
	return n == 1
}

const groupMessageCount = `
select sum(messages) from real_users
where gid = ?`

func (m *Model) GroupMessageCount(g Group) (int, error) {
	var c int
	err := m.db.Get(&c, groupMessageCount, g.GID)
	return c, err
}

const groupAverageMessageCount = `
select avg(messages) from real_users
where gid = ?`

func (m *Model) AverageMessageCount(g Group) (float64, error) {
	var c float64
	err := m.db.Get(&c, groupAverageMessageCount, g.GID)
	return c, err
}
