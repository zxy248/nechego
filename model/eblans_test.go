package model

import (
	"errors"
	"testing"
)

func FuzzEblans(f *testing.F) {
	db := testingDB()
	defer db.DropTables()
	e := &Eblans{db}

	testcases := []struct {
		gid int64
		uid int64
	}{{1, 123}, {2, 234}, {-1234, 12341243}, {-78123478, 829345}}
	for _, tc := range testcases {
		f.Add(tc.gid, tc.uid)
	}
	f.Fuzz(func(t *testing.T, gid, uid int64) {
		_, err := e.Get(gid)
		if !errors.Is(err, ErrNoEblan) {
			t.Errorf("err = %v, want %v", err, ErrNoEblan)
		}

		if err := e.Insert(gid, uid); err != nil {
			t.Fatal(err)
		}

		id, err := e.Get(gid)
		if err != nil {
			t.Fatal(err)
		}
		if id != uid {
			t.Errorf("id = %v, want %v", id, uid)
		}

		if err := e.Delete(gid); err != nil {
			t.Fatal(err)
		}
	})
}
