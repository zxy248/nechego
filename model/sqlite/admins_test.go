package sqlite

import (
	"fmt"
	"testing"
)

func TestAdminsList(t *testing.T) {
	db := testingDB()
	defer db.DropTables()
	a := &Admins{db}
	gid := int64(0)

	testcases := []int64{1, 234, 12341243, 829345, 100000000}
	for i, tc := range testcases {
		t.Run(fmt.Sprintf("testcase %v", tc), func(t *testing.T) {
			l, err := a.List(gid)
			if err != nil {
				t.Fatal(err)
			}
			if len(l) != i {
				t.Errorf("len(l) = %v, want %v", len(l), i)
			}

			if err := a.Insert(tc); err != nil {
				t.Fatal(err)
			}

			l, err = a.List(gid)
			if err != nil {
				t.Fatal(err)
			}
			if len(l) != i+1 {
				t.Errorf("len(l) = %v, want %v", len(l), i+1)
			}
			if l[len(l)-1] != tc {
				t.Errorf("l[last] = %v, want %v", l[len(l)-1], tc)
			}
		})
	}
}

func FuzzAdminsAuthorize(f *testing.F) {
	db := testingDB()
	defer db.DropTables()
	a := &Admins{db}
	gid := int64(0)

	testcases := []int64{1, 234, 7981234}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, id int64) {
		authorized, err := a.Authorize(gid, id)
		if err != nil {
			t.Fatal(err)
		}
		if authorized {
			t.Error("authorized before insertion")
		}

		if err := a.Insert(id); err != nil {
			t.Fatal(err)
		}

		authorized, err = a.Authorize(gid, id)
		if err != nil {
			t.Fatal(err)
		}
		if !authorized {
			t.Errorf("not authorized after insertion")
		}

		if err := a.Delete(id); err != nil {
			t.Fatal(err)
		}

		authorized, err = a.Authorize(gid, id)
		if err != nil {
			t.Fatal(err)
		}
		if authorized {
			t.Errorf("authorized after deletion")
		}
	})
}
