package model

import (
	"database/sql"
	"errors"
)

var ErrNoPair = errors.New("no pair")
var ErrNoEblan = errors.New("no eblan")

type DB struct {
	*sql.DB
}

const createTablesQuery = `
create table if not exists users (
    id integer,
    gid integer not null,
    uid integer not null,
    primary key (id autoincrement)
);

create table if not exists pairs (
    id integer not null,
    gid integer not null,
    uidx integer not null,
    uidy integer not null,
    added text not null,
    primary key (id autoincrement)
);

create table if not exists eblans (
    id integer not null,
    gid integer not null,
    uid integer not null,
    added text not null,
    primary key (id autoincrement)
);

create table if not exists whitelist (
    id integer not null,
    gid integer not null,
    primary key (id autoincrement)
);

create table if not exists admins (
    id integer not null,
    uid integer not null,
    primary key (id autoincrement)
);

create table if not exists bans (
    id integer not null,
    uid integer not null,
    primary key (id autoincrement)
);

create table if not exists status (
    id integer not null,
    gid integer not null,
    primary key (id autoincrement)
);
`

func (db *DB) Setup() error {
	_, err := db.Exec(createTablesQuery)
	if err != nil {
		return err
	}
	return nil
}

const teardownQuery = `
drop table users;
drop table pairs;
drop table eblans;
drop table whitelist;
drop table admins;
drop table bans;
drop table status;
`

func (db *DB) Teardown() error {
	_, err := db.Exec(teardownQuery)
	if err != nil {
		return err
	}
	return nil
}

type Users struct {
	DB *DB
}

const insertUserQuery = `
insert into users (uid, gid) values (?, ?)
`

// Insert inserts a user to the users table.
func (u *Users) Insert(gid, uid int64) error {
	_, err := u.DB.Exec(insertUserQuery, uid, gid)
	if err != nil {
		return err
	}
	return nil
}

const deleteUserQuery = `
delete from users where gid = ? and uid = ?
`

func (u *Users) Delete(gid, uid int64) error {
	_, err := u.DB.Exec(deleteUserQuery, gid, uid)
	if err != nil {
		return err
	}
	return nil
}

const listUsersQuery = `
select uid from users where gid = ?
`

// List gets returns all users from the users table.
func (u *Users) List(gid int64) ([]int64, error) {
	rows, err := u.DB.Query(listUsersQuery, gid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

const existsUserQuery = `
select 1 from users where gid = ? and uid = ?
`

func (u *Users) Exists(gid, uid int64) (bool, error) {
	var i int
	if err := u.DB.QueryRow(existsUserQuery, gid, uid).Scan(&i); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return i == 1, nil
}

const randomUserQuery = `
select uid from users where gid = ? order by random() limit 1
`

func (u *Users) Random(gid int64) (int64, error) {
	var uid int64
	if err := u.DB.QueryRow(randomUserQuery, gid).Scan(&uid); err != nil {
		return 0, err
	}
	return uid, nil
}

type Pairs struct {
	DB *DB
}

const insertPairQuery = `
insert into pairs (gid, uidx, uidy, added)
values (?, ?, ?, datetime('now', 'localtime'))
`

// Insert inserts the pair to the pairs table.
func (p *Pairs) Insert(gid int64, uidx, uidy int64) error {
	_, err := p.DB.Exec(insertPairQuery, gid, uidx, uidy)
	if err != nil {
		return err
	}
	return nil
}

const getPairQuery = `
select uidx, uidy from pairs
where gid = ? and added > date('now', 'localtime')
order by added desc
limit 1
`

// Get gets the pair from the pairs table.
func (p *Pairs) Get(gid int64) (uidx, uidy int64, err error) {
	if err := p.DB.QueryRow(getPairQuery, gid).Scan(&uidx, &uidy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, ErrNoPair
		}
		return 0, 0, err
	}
	return uidx, uidy, nil
}

type Eblans struct {
	DB *DB
}

const insertEblanQuery = `
insert into eblans (gid, uid, added)
values (?, ?, datetime('now', 'localtime'))
`

// Insert inserts the user to the eblans table.
func (e *Eblans) Insert(gid, uid int64) error {
	_, err := e.DB.Exec(insertEblanQuery, gid, uid)
	if err != nil {
		return err
	}
	return nil
}

const getEblanQuery = `
select uid from eblans
where gid = ? and added > date('now', 'localtime')
order by added desc
limit 1
`

// Get gets the user from the eblans table.
func (e *Eblans) Get(gid int64) (int64, error) {
	var uid int64
	if err := e.DB.QueryRow(getEblanQuery, gid).Scan(&uid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoEblan
		}
		return 0, err
	}
	return uid, nil
}

type Whitelist struct {
	DB *DB
}

const insertWhitelistQuery = `
insert into whitelist (gid) values (?)
`

// Insert inserts a group ID to the whitelist table.
func (w *Whitelist) Insert(gid int64) error {
	_, err := w.DB.Exec(insertWhitelistQuery, gid)
	if err != nil {
		return err
	}
	return nil
}

const allowWhitelistQuery = `
select 1 from whitelist where gid = ?
`

// Allow returns true if the given group ID exists in the whitelist table.
func (w *Whitelist) Allow(gid int64) (bool, error) {
	var i int
	if err := w.DB.QueryRow(allowWhitelistQuery, gid).Scan(&i); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return i == 1, nil
}

type Admins struct {
	DB *DB
}

const insertAdminsQuery = `
insert into admins (uid) values (?)
`

// Insert inserts a user ID to the admins table.
func (a *Admins) Insert(uid int64) error {
	_, err := a.DB.Exec(insertAdminsQuery, uid)
	if err != nil {
		return err
	}
	return nil
}

const listAdminsQuery = `
select uid from admins
`

// List returns the list of admins.
func (a *Admins) List() ([]int64, error) {
	rows, err := a.DB.Query(listAdminsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

const allowAdminsQuery = `
select 1 from admins where uid = ?
`

func (a *Admins) Allow(uid int64) (bool, error) {
	var i int
	if err := a.DB.QueryRow(allowAdminsQuery, uid).Scan(&i); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return i == 1, nil
}

type Bans struct {
	DB *DB
}

const insertBansQuery = `
insert into bans (uid) values (?)
`

func (b *Bans) Ban(uid int64) error {
	_, err := b.DB.Exec(insertBansQuery, uid)
	if err != nil {
		return err
	}
	return nil
}

const deleteBansQuery = `
delete from bans where uid = ?
`

func (b *Bans) Unban(uid int64) error {
	_, err := b.DB.Exec(deleteBansQuery, uid)
	if err != nil {
		return err
	}
	return nil
}

const listBansQuery = `
select uid from bans
`

// List returns a list of banned users.
func (b *Bans) List() ([]int64, error) {
	rows, err := b.DB.Query(listBansQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

const bannedBansQuery = `
select 1 from bans where uid = ?
`

func (b *Bans) Banned(uid int64) (bool, error) {
	var i int
	if err := b.DB.QueryRow(bannedBansQuery, uid).Scan(&i); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return i == 1, nil
}

type Status struct {
	DB *DB
}

const enableStatusQuery = `
delete from status where gid = ?
`

func (s *Status) Enable(gid int64) error {
	_, err := s.DB.Exec(enableStatusQuery, gid)
	if err != nil {
		return err
	}
	return nil
}

const activeStatusQuery = `
select 1 from status where gid = ?
`

func (s *Status) Active(gid int64) (bool, error) {
	var i int
	if err := s.DB.QueryRow(activeStatusQuery, gid).Scan(&i); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, err
	}
	return !(i == 1), nil
}

const disableStatusQuery = `
insert into status (gid) values (?)
`

func (s *Status) Disable(gid int64) error {
	_, err := s.DB.Exec(disableStatusQuery, gid)
	if err != nil {
		return err
	}
	return nil
}
