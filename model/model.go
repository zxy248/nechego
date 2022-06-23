package model

import (
	"database/sql"
)

type DB struct {
	*sql.DB
}

const createTableUsersQuery = `
create table if not exists users (
    id integer,
    gid integer not null,
    uid integer not null,
    primary key (id autoincrement)
)
`

const createTablePairsQuery = `
create table if not exists pairs (
    id integer not null,
    gid integer not null,
    uidx integer not null,
    uidy integer not null,
    added text not null,
    primary key (id autoincrement)
)
`

const createTableEblansQuery = `
create table if not exists eblans (
    id integer not null,
    gid integer not null,
    uid integer not null,
    added text not null,
    primary key (id autoincrement)
)
`

const createTableWhitelistQuery = `
create table if not exists whitelist (
    id integer not null,
    gid integer not null,
    primary key (id autoincrement)
)
`

const createTableAdminsQuery = `
create table if not exists admins (
    id integer not null,
    uid integer not null,
    primary key (id autoincrement)
)
`

const createTableBansQuery = `
create table if not exists bans (
    id integer not null,
    uid integer not null,
    primary key (id autoincrement)
)
`

const createTableStatusQuery = `
create table if not exists status (
    id integer not null,
    gid integer not null,
    primary key (id autoincrement)
)
`

const createTableForbidQuery = `
create table if not exists forbid (
    id integer not null,
    gid integer not null,
    command integer not null,
    primary key (id autoincrement)
)`

const createTableAdminQuery = `
create table if not exists admin (
    id integer not null,
    gid integer not null,
    uid integer not null,
    added text not null,
    primary key (id autoincrement)
)
`

// CreateTables creates the necessary tables.
func (db *DB) CreateTables() error {
	queries := []string{
		createTableUsersQuery,
		createTablePairsQuery,
		createTableEblansQuery,
		createTableWhitelistQuery,
		createTableAdminsQuery,
		createTableBansQuery,
		createTableStatusQuery,
		createTableForbidQuery,
		createTableAdminQuery,
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

const (
	dropTableUsersQuery     = "drop table users"
	dropTablePairsQuery     = "drop table pairs"
	dropTableEblansQuery    = "drop table eblans"
	dropTableWhitelistQuery = "drop table whitelist"
	dropTableAdminsQuery    = "drop table admins"
	dropTableBansQuery      = "drop table bans"
	dropTableStatusQuery    = "drop table status"
	dropTableForbidQuery    = "drop table forbid"
	dropTableAdminQuery     = "drop table admin"
)

// DropTables deletes all tables from the database.
func (db *DB) DropTables() error {
	queries := []string{
		dropTableUsersQuery,
		dropTablePairsQuery,
		dropTableEblansQuery,
		dropTableWhitelistQuery,
		dropTableAdminsQuery,
		dropTableBansQuery,
		dropTableStatusQuery,
		dropTableForbidQuery,
		dropTableAdminQuery,
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}
