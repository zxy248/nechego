package model

import "testing"

func FuzzBansBanned(f *testing.F) {
	db := testingDB()
	defer db.DropTables()
	b := &Bans{db}

	testcases := []int64{1, 234, 12341243, 829345, 100000000}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, id int64) {
		banned, err := b.Banned(id)
		if err != nil {
			t.Fatal(err)
		}
		if banned {
			t.Error("banned before ban")
		}

		if err := b.Ban(id); err != nil {
			t.Fatal(err)
		}

		banned, err = b.Banned(id)
		if err != nil {
			t.Fatal(err)
		}
		if !banned {
			t.Error("not banned after ban")
		}

		if err := b.Unban(id); err != nil {
			t.Fatal(err)
		}

		banned, err = b.Banned(id)
		if err != nil {
			t.Fatal(err)
		}
		if banned {
			t.Error("banned after unban")
		}
	})
}

func FuzzBansList(f *testing.F) {
	db := testingDB()
	defer db.DropTables()
	b := &Bans{db}

	i := 0
	testcases := []int64{1, 234, 12341243, 829345, 100000000}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, id int64) {
		l, err := b.List()
		if err != nil {
			t.Fatal(err)
		}
		if len(l) != i {
			t.Errorf("len(l) = %v, want %v", len(l), i)
		}

		if err := b.Ban(id); err != nil {
			t.Fatal(err)
		}
		i++

		l, err = b.List()
		if err != nil {
			t.Fatal(err)
		}
		if len(l) != i {
			t.Errorf("len(l) = %v, want %v", len(l), i)
		}
		if l[len(l)-1] != id {
			t.Errorf("l[last] = %v, want %v", l[len(l)-1], id)
		}
	})
}
