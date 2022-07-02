package sqlite

import "testing"

func TestStatus(t *testing.T) {
	db := testingDB()
	defer db.DropTables()
	s := &Status{db}

	gid := int64(1234)

	on, err := s.Active(gid)
	if err != nil {
		t.Fatal(err)
	}
	if !on {
		t.Errorf("off, want on")
	}

	if err := s.Disable(gid); err != nil {
		t.Fatal(err)
	}

	on, err = s.Active(gid)
	if err != nil {
		t.Fatal(err)
	}
	if on {
		t.Errorf("on, want off")
	}

	if err := s.Enable(gid); err != nil {
		t.Fatal(err)
	}

	on, err = s.Active(gid)
	if err != nil {
		t.Fatal(err)
	}
	if !on {
		t.Errorf("off, want on")
	}
}
