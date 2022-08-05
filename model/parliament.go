package model

import (
	"errors"
)

func (m *Model) Parliament(g Group, n int) ([]User, error) {
	tx := m.db.MustBegin()
	defer tx.Rollback()
	users := []User{}
	if err := tx.Select(&users, selectEvents, parliamentMemberEvent, g.GID); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		if err := tx.Select(&users, randomUsers, g.GID, n); err != nil {
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

	var member bool
	if err := tx.Get(&member, existsUserEventToday, parliamentMemberEvent, g.GID, u.ID); err != nil {
		return 0, err
	}
	if !member {
		return 0, ErrNotParliamentMember
	}

	var voted bool
	if err := tx.Get(&voted, existsUserEventToday, impeachmentEvent, g.GID, u.ID); err != nil {
		return 0, err
	}
	if voted {
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
