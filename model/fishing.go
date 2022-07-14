package model

import (
	"nechego/fishing"
)

type Catch struct {
	ID     int
	UserID int `db:"user_id"`
	Fish   fishing.Fish
	Sold   bool
	Frozen bool
}

func MakeCatch(u User, f fishing.Fish) Catch {
	return Catch{
		UserID: u.ID,
		Fish:   f,
	}
}

const insertFish = `
insert into fishing (user_id, fish, sold, frozen)
values (?, ?, ?, ?)`

func (m *Model) InsertFish(c Catch) {
	m.db.MustExec(insertFish, c.UserID, c.Fish, c.Sold, c.Frozen)
}

const selectFish = `
select id, user_id, fish, sold, frozen from fishing
where user_id = ? and sold = 0`

func (m *Model) SelectFish(u User) ([]Catch, error) {
	c := []Catch{}
	if err := m.db.Select(&c, selectFish, u.ID); err != nil {
		return nil, err
	}
	return c, nil
}

const unsoldFish = selectFish + `
and frozen = 0`

const sellFish = `
update fishing set sold = 1
where id = ?`

func (m *Model) SellFish(u User) ([]Catch, error) {
	tx := m.db.MustBegin()
	defer tx.Rollback()

	catch := []Catch{}
	if err := tx.Select(&catch, unsoldFish, u.ID); err != nil {
		return nil, err
	}
	for _, c := range catch {
		tx.MustExec(sellFish, c.ID)
	}
	
	return catch, tx.Commit()
}

const freezeFish = `
update fishing set frozen = 1
where user_id = ? and sold = 0 and frozen = 0`

func (m *Model) FreezeFish(u User) {
	m.db.MustExec(freezeFish, u.ID)
}

const unfreezeFish = `
update fishing set frozen = 0
where user_id = ? and sold = 0 and frozen = 1`

func (m *Model) UnfreezeFish(u User) {
	m.db.MustExec(unfreezeFish, u.ID)
}
