package model

import (
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestUsers(t *testing.T) {
	db := testingDB(t)
	defer db.Teardown()

	u := &Users{db}
	ids, err := u.List(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 0 {
		t.Errorf("len(ids) == %v, %v expected\n", len(ids), 0)
	}

	if err := u.Insert(1, 123); err != nil {
		t.Fatal(err)
	}
	ids, err = u.List(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 {
		t.Errorf("len(ids) == %v, %v expected\n", len(ids), 0)
	}
	if ids[0] != 123 {
		t.Errorf("ids[0] == %v, %v expected\n", ids[0], 123)
	}
}

func TestPairs(t *testing.T) {
	db := testingDB(t)
	defer db.Teardown()

	p := &Pairs{db}
	_, _, err := p.Get(1)
	if !errors.Is(err, ErrNoPair) {
		t.Errorf("err == %v, %v expected", err, ErrNoPair)
	}

	if err := p.Insert(1, 123, 124); err != nil {
		t.Fatal(err)
	}
	x, y, err := p.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if x != 123 || y != 124 {
		t.Errorf("x, y == %v, %v; %v %v expected", x, y, 123, 124)
	}
}

func TestEblans(t *testing.T) {
	db := testingDB(t)
	defer db.Teardown()

	e := &Eblans{db}
	_, err := e.Get(1)
	if !errors.Is(err, ErrNoEblan) {
		t.Errorf("err == %v, %v expected", err, ErrNoEblan)
	}

	if err := e.Insert(1, 123); err != nil {
		t.Fatal(err)
	}
	uid, err := e.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if uid != 123 {
		t.Errorf("uid == %v, %v expected", uid, 123)
	}
}

func TestWhitelist(t *testing.T) {
	db := testingDB(t)
	defer db.Teardown()

	w := &Whitelist{db}
	ok, err := w.Allow(1)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Errorf("ok == %v, %v expected", ok, false)
	}

	if err := w.Insert(1); err != nil {
		t.Fatal(err)
	}

	ok, err = w.Allow(1)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("ok == %v, %v expected", ok, true)
	}
}

func TestAdmins(t *testing.T) {
	db := testingDB(t)
	defer db.Teardown()

	a := &Admins{db}

	ok, err := a.Allow(1)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Errorf("ok == %v, %v expected", ok, false)
	}

	l, err := a.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 0 {
		t.Errorf("len(l) == %v, %v expected", len(l), 0)
	}

	if err := a.Insert(1); err != nil {
		t.Fatal(err)
	}

	ok, err = a.Allow(1)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("ok == %v, %v expected", ok, true)
	}

	l, err = a.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 1 {
		t.Errorf("len(l) == %v, %v expected", len(l), 1)
	}
	if l[0] != 1 {
		t.Errorf("l[0] == %v, %v expected", l[0], 1)
	}
}

func TestBans(t *testing.T) {
	db := testingDB(t)
	defer db.Teardown()

	b := &Bans{db}

	ok, err := b.Banned(1)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Errorf("ok == %v, %v expected", ok, false)
	}

	l, err := b.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 0 {
		t.Errorf("len(l) == %v, %v expected", len(l), 0)
	}

	if err := b.Ban(1); err != nil {
		t.Fatal(err)
	}

	ok, err = b.Banned(1)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("ok == %v, %v expected", ok, true)
	}

	l, err = b.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 1 {
		t.Errorf("len(l) == %v, %v expected", len(l), 1)
	}
	if l[0] != 1 {
		t.Errorf("l[0] == %v, %v expected", l[0], 1)
	}
}

func TestStatus(t *testing.T) {
	db := testingDB(t)
	defer db.Teardown()

	s := &Status{db}

	ok, err := s.Active(1)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("ok == %v, %v expected", ok, true)
	}

	if err := s.Disable(1); err != nil {
		t.Fatal(err)
	}

	ok, err = s.Active(1)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Errorf("ok == %v, %v expected", ok, false)
	}

	if err := s.Enable(1); err != nil {
		t.Fatal(err)
	}

	ok, err = s.Active(1)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("ok == %v, %v expected", ok, true)
	}
}

func testingDB(t *testing.T) *DB {
	sqlDB, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	modelDB := &DB{sqlDB}
	if err := modelDB.Setup(); err != nil {
		t.Fatal(err)
	}
	return modelDB
}
