package sqlite

import (
	"errors"
	"nechego/model"
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
	if !errors.Is(err, model.ErrNoPair) {
		t.Errorf("err = %v, want %v", err, model.ErrNoPair)
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
