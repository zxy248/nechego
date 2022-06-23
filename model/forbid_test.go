package model

import (
	"nechego/input"
	"testing"
)

func FuzzForbidForbidden(f *testing.F) {
	db := testingDB()
	defer db.DropTables()
	m := &Forbid{db}

	testcases := []struct {
		gid     int64
		command int
	}{
		{1, 5}, {719823479, 404}, {123456789, 300},
	}
	for _, tc := range testcases {
		f.Add(tc.gid, tc.command)
	}
	f.Fuzz(func(t *testing.T, gid int64, c int) {
		command := input.Command(c)

		forbidden, err := m.Forbidden(gid, command)
		if err != nil {
			t.Fatal(err)
		}
		if forbidden {
			t.Errorf("forbidden before forbid")
		}

		if err := m.Forbid(gid, command); err != nil {
			t.Fatal(err)
		}

		forbidden, err = m.Forbidden(gid, command)
		if err != nil {
			t.Fatal(err)
		}
		if !forbidden {
			t.Errorf("not forbidden after forbid")
		}

		if err := m.Permit(gid, command); err != nil {
			t.Fatal(err)
		}

		forbidden, err = m.Forbidden(gid, command)
		if err != nil {
			t.Fatal(err)
		}
		if forbidden {
			t.Errorf("forbidden after permit")
		}
	})
}

func FuzzForbidList(f *testing.F) {
	db := testingDB()
	defer db.DropTables()
	m := &Forbid{db}

	i := 0
	testcases := []struct {
		gid     int64
		command int
	}{
		{1, 5}, {719823479, 404}, {123456789, 300},
	}
	for _, tc := range testcases {
		f.Add(tc.gid, tc.command)
	}
	f.Fuzz(func(t *testing.T, gid int64, c int) {
		command := input.Command(c)

		l, err := m.List(gid)
		if err != nil {
			t.Fatal(err)
		}
		if len(l) != i {
			t.Errorf("len(l) = %v, want %v", len(l), i)
		}

		if err := m.Forbid(gid, command); err != nil {
			t.Fatal(err)
		}
		i++

		l, err = m.List(gid)
		if err != nil {
			t.Fatal(err)
		}
		if len(l) != i {
			t.Errorf("len(l) = %v, want %v", len(l), i)
		}
		if l[len(l)-1] != command {
			t.Errorf("l[last] = %v, want %v", l[len(l)-1], command)
		}

		if err := m.Permit(gid, command); err != nil {
			t.Fatal(err)
		}
		i--

		l, err = m.List(gid)
		if err != nil {
			t.Fatal(err)
		}
		if len(l) != i {
			t.Errorf("len(l) = %v, want %v", len(l), i)
		}
	})
}
