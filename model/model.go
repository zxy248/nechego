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

func New(db *sqlx.DB) *Model {
	d := &DB{db}
	d.CreateTables()
	return &Model{d}
}

// don't forget to update views when updating tables
const schema = `
create table if not exists users (
    id integer primary key autoincrement,
    gid integer not null references groups (gid) on delete cascade,
    uid integer not null,
    energy integer not null default 0 check (energy >= 0),
    balance integer not null default 0 check (balance >= 0),
    account integer not null default 0 check (account >= 0),
    debt integer not null default 0 check (debt >= 0),
    debt_limit integer not null default 0 check (debt_limit >= 0),
    admin integer not null default 0,
    banned integer not null default 0,
    messages integer not null default 0,
    fisher integer not null default 0,
    fishes integer not null default 0,
    active integer not null default 1,
    unique (gid, uid)
);

create table if not exists groups (
    gid integer primary key,
    whitelisted integer not null default 0,
    status integer not null default 1
);

create table if not exists daily_pairs (
    id integer primary key autoincrement,
    gid integer not null references groups (gid) on delete cascade,
    user_id_x integer not null references users (id) on delete cascade,
    user_id_y integer not null references users (id) on delete cascade,
    added datetime not null,
    check (user_id_x != user_id_y)
);

create table if not exists daily_eblans (
    id integer primary key autoincrement,
    gid integer not null references groups (gid) on delete cascade,
    user_id integer not null references users (id) on delete cascade,
    added datetime not null
);

create table if not exists daily_admins (
    id integer primary key autoincrement,
    gid integer not null references groups (gid) on delete cascade,
    user_id integer not null references users (id) on delete cascade,
    added datetime not null
);

create table if not exists forbidden_commands (
    id integer primary key autoincrement,
    gid integer not null references groups (gid) on delete cascade,
    command integer not null,
    unique (gid, command)
);

create view if not exists real_users as
select id, gid, uid, energy, balance, account, debt, debt_limit, admin
or exists(select 1 from daily_admins
    where daily_admins.user_id = users.id
    and added >= date('now', 'localtime'))
as admin,
banned, messages, fisher, fishes
from users
where active = 1;

create table if not exists events (
    id integer primary key autoincrement,
    gid integer not null references groups (gid) on delete cascade,
    user_id integer not null references users (id) on delete cascade,
    event integer not null,
    happen datetime not null
);

create table if not exists fishing (
    id integer primary key autoincrement,
    user_id integer not null references users (id) on delete cascade,
    fish text not null,
    sold integer not null,
    frozen integer not null
);

create table if not exists pets (
    id integer primary key autoincrement,
    user_id integer not null references users (id) on delete cascade,
    name text not null,
    species integer not null,
    gender integer not null,
    birth datetime not null,
    unique (user_id)
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
drop view real_users;
drop table events;
drop table fishing;
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
