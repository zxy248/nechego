package model

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

type Model struct {
	db *DB
}

func NewModel(db *sqlx.DB) *Model {
	d := &DB{db}
	d.CreateTables()
	return &Model{d}
}

const schema = `
create table if not exists users (
    id integer primary key autoincrement,
    gid integer references groups (gid) on delete cascade,
    uid integer not null,
    energy integer not null check (energy >= 0),
    balance integer not null check (balance >= 0),
    admin integer not null,
    banned integer not null,
    messages integer not null,
    unique (gid, uid)
);

create table if not exists groups (
    gid integer primary key,
    whitelisted integer not null,
    status integer not null
);

create table if not exists daily_pairs (
    id integer primary key autoincrement,
    gid integer references groups (gid) on delete cascade,
    user_id_x integer references users (id) on delete cascade,
    user_id_y integer references users (id) on delete cascade,
    added datetime not null,
    check (user_id_x != user_id_y)
);

create table if not exists daily_eblans (
    id integer primary key autoincrement,
    gid integer references groups (gid) on delete cascade,
    user_id integer references users (id) on delete cascade,
    added datetime not null
);

create table if not exists daily_admins (
    id integer primary key autoincrement,
    gid integer references groups (gid) on delete cascade,
    user_id integer references users (id) on delete cascade,
    added datetime not null
);

create table if not exists forbidden_commands (
    id integer primary key autoincrement,
    gid integer references groups (gid) on delete cascade,
    command integer not null,
    unique (gid, command)
);
`

// CreateTables creates the necessary tables.
func (db *DB) CreateTables() {
	db.MustExec(schema)
}

const drop = `
drop table users;
drop table groups;
drop table daily_pairs;
drop table daily_eblans;
drop table daily_admins;
drop table forbidden_commands;
`

// DropTables deletes all tables from the database.
func (db *DB) DropTables() {
	db.MustExec(drop)
}

func concat(elems ...string) string {
	return strings.Join(elems, " ")
}

func failOn(err error) {
	if err != nil {
		panic(fmt.Errorf("unexpected model error: %v", err))
	}
}
