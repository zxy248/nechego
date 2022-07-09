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
update groups set status = ?
where gid = ?`

func (m *Model) EnableGroup(g Group) {
	m.db.MustExec(enableGroup, true, g.GID)
}

func (m *Model) DisableGroup(g Group) {
	m.db.MustExec(enableGroup, false, g.GID)
}
