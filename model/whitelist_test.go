package model

import "testing"

func TestWhitelist(t *testing.T) {
	db := testingDB()
	defer db.DropTables()
	w := &Whitelist{db}

	gid := int64(1234)

	ok, err := w.Allow(gid)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Errorf("ok, want !ok")
	}

	if err := w.Insert(gid); err != nil {
		t.Fatal(err)
	}

	ok, err = w.Allow(gid)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("!ok, want ok")
	}
}
