package model

import (
	"database/sql"
	"errors"
)

const selectParliament = `
select u.* from real_users as u
inner join events as e on u.id = e.user_id
where e.event = ?
and e.gid = ?
and e.happen >= date('now', 'localtime')`

const insertEvent = `
insert into events (gid, user_id, event, happen)
values (?, ?, ?, datetime('now', 'localtime'))`

func (m *Model) Parliament(g Group, n int) ([]User, error) {
	tx := m.db.MustBegin()
	defer tx.Rollback()
	users := []User{}
	err := tx.Select(&users, selectParliament, parliamentMemberEvent, g.GID)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		err := tx.Select(&users, randomUsers, g.GID, n)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			tx.MustExec(insertEvent, g.GID, u.ID, parliamentMemberEvent)
		}
	}
	return users, tx.Commit()
}

const countEvents = `
select count(1) from events
where event = ?
and gid = ?
and happen >= date('now', 'localtime')`

const impeachedToday = `
select 1 from events
where event = ?
and gid = ?
and user_id = ?
and happen >= date('now', 'localtime')`

const cancelAdmin = `
delete from daily_admins
where gid = ?
and added >= date('now', 'localtime')`

var (
	ErrNotParliamentMember = errors.New("not a parliament member")
	ErrAlreadyVoted        = errors.New("already voted")
	ErrAlreadyImpeached    = errors.New("already impeached")
)

func (m *Model) Impeachment(g Group, u User, threshold int) (votes int, err error) {
	tx := m.db.MustBegin()
	defer tx.Rollback()

	parliament := []User{}
	if err := tx.Select(&parliament, selectParliament, parliamentMemberEvent, g.GID); err != nil {
		return 0, err
	}
	e := false
	for _, p := range parliament {
		if p.ID == u.ID {
			e = true
			break
		}
	}
	if !e {
		return 0, ErrNotParliamentMember
	}

	var test int
	err = tx.Get(&test, impeachedToday, impeachmentEvent, g.GID, u.ID)
	if !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			return 0, err
		}
		return 0, ErrAlreadyVoted
	}

	tx.MustExec(insertEvent, g.GID, u.ID, impeachmentEvent)

	var c int
	if err := tx.Get(&c, countEvents, impeachmentEvent, g.GID); err != nil {
		return 0, err
	}
	if c == threshold {
		tx.MustExec(cancelAdmin, g.GID)
	}
	if c > threshold {
		return 0, ErrAlreadyImpeached
	}
	return c, tx.Commit()
}
