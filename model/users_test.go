package model

import "testing"

func TestUsers(t *testing.T) {
	db := testingDB()
	defer db.DropTables()
	u := &Users{db}

	gid, uid := int64(1234), int64(7981234)

	ids, err := u.List(gid)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 0 {
		t.Errorf("len(ids) = %v, want empty", len(ids))
	}

	if err := u.Insert(gid, uid); err != nil {
		t.Fatal(err)
	}

	ids, err = u.List(gid)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 {
		t.Errorf("len(ids) = %v, want %v\n", len(ids), 1)
	}
	if ids[0] != uid {
		t.Errorf("ids[0] = %v, want %v\n", ids[0], uid)
	}
}

func TestNRandomUsers(t *testing.T) {
	db := testingDB()
	defer db.DropTables()
	u := &Users{db}
	const (
		gid = 0
		n   = 6
		j   = 4
	)
	i := int64(0)

	l, err := u.NRandom(gid, n)
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 0 {
		t.Error("l is not empty")
	}

	for ; i < j; i++ {
		if err := u.Insert(gid, i); err != nil {
			t.Fatal(err)
		}
	}
	l, err = u.NRandom(gid, n)
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != j {
		t.Errorf("len(l) == %v, want %v", len(l), j)
	}

	for ; i < 2*j; i++ {
		if err := u.Insert(gid, i); err != nil {
			t.Fatal(err)
		}
	}
	l, err = u.NRandom(gid, n)
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != n {
		t.Errorf("len(l) == %v, want %v", len(l), n)
	}
}
