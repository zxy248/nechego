package model

import (
	"errors"
	"testing"
)

func TestPairs(t *testing.T) {
	db := testingDB()
	defer db.DropTables()
	p := &Pairs{db}

	gid := int64(1)
	uidx := int64(123)
	uidy := int64(234)

	_, _, err := p.Get(gid)
	if !errors.Is(err, ErrNoPair) {
		t.Errorf("err = %v, want %v", err, ErrNoPair)
	}

	if err := p.Insert(gid, uidx, uidy); err != nil {
		t.Fatal(err)
	}

	x, y, err := p.Get(gid)
	if err != nil {
		t.Fatal(err)
	}
	if x != uidx || y != uidy {
		t.Errorf("x, y = %v, %v; want %v, %v", x, y, uidx, uidy)
	}
}
