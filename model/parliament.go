package model

import (
	"database/sql"
	"errors"
)

func (m *Model) Parliament(g Group, n int) ([]User, error) {
	tx := m.db.MustBegin()
	defer tx.Rollback()
	users := []User{}
	err := tx.Select(&users, selectEvents, parliamentMemberEvent, g.GID)
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
	if err := tx.Select(&parliament, selectEvents, parliamentMemberEvent, g.GID); err != nil {
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
	err = tx.Get(&test, existsUserEventToday, impeachmentEvent, g.GID, u.ID)
	if !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			return 0, err
		}
		return 0, ErrAlreadyVoted
	}

	tx.MustExec(insertEvent, g.GID, u.ID, impeachmentEvent)

	var c int
	if err := tx.Get(&c, countEventsToday, impeachmentEvent, g.GID); err != nil {
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
