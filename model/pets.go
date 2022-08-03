package model

import (
	"database/sql"
	"errors"
	"nechego/pets"
)

const insertPet = `
insert into pets (user_id, name, species, gender, birth)
values (?, ?, ?, ?, ?)`

func (m *Model) InsertPet(u User, p *pets.Pet) {
	m.db.MustExec(insertPet, u.ID, p.Name, p.Species, p.Gender, p.Birth)
}

const getPet = `
select name, species, gender, birth from pets
where user_id = ?`

var ErrNoPet = errors.New("no pet")

func (m *Model) GetPet(u User) (*pets.Pet, error) {
	var p pets.Pet
	if err := m.db.Get(&p, getPet, u.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoPet
		}
		return nil, err
	}
	return &p, nil
}

const deletePet = `
delete from pets
where user_id = ?`

func (m *Model) DeletePet(u User) {
	m.db.MustExec(deletePet, u.ID)
}

const namePet = `
update pets set name = ?
where user_id = ?`

func (m *Model) NamePet(u User, s string) {
	m.db.MustExec(namePet, s, u.ID)
}

const hasPet = `
select exists(
select 1 from pets
where user_id = ?
limit 1)`

func (m *Model) HasPet(u User) (bool, error) {
	var flag bool
	err := m.db.Get(&flag, hasPet, u.ID)
	return flag, err
}
