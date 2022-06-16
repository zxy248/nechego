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
