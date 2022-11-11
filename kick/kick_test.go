package kick

import (
	"testing"
	"time"
)

func getTimeoutFunc(flags []bool, i int) func() {
	return func() {
		flags[i] = true
	}
}

func TestVote(t *testing.T) {
	duration := time.Millisecond * 100
	delta := time.Millisecond * 10
	votes := 3
	k := New(duration, votes)

	groups := []int64{42, 75}
	users := []int64{31, 98, 33, 94, 34}
	targets := []int64{83, 62, 22}
	flags := []bool{false, false}
	timeouts := []func(){getTimeoutFunc(flags, 0), getTimeoutFunc(flags, 1)}

	s := k.Vote(groups[0], users[0], 0, timeouts[0])
	wantNoUser := Status{NoUser, 0, 0}
	if s != wantNoUser {
		t.Errorf("s == %v, want %v", s, wantNoUser)
	}

	s = k.Vote(groups[1], users[0], targets[0], timeouts[1])
	wantInit1 := Status{Init, targets[0], votes - 1}
	if s != wantInit1 {
		t.Errorf("s == %v, want %v", s, wantInit1)
	}

	s = k.Vote(groups[0], users[0], targets[0], timeouts[0])
	wantInit := Status{Init, targets[0], votes - 1}
	if s != wantInit {
		t.Errorf("s == %v, want %v", s, wantInit)
	}

	s = k.Vote(groups[0], users[0], targets[0], timeouts[0])
	wantDuplicate := Status{Duplicate, targets[0], votes - 1}
	if s != wantDuplicate {
		t.Errorf("s == %v, want %v", s, wantDuplicate)
	}

	s = k.Vote(groups[0], users[0], targets[1], timeouts[0])
	wantConflict := Status{Conflict, targets[0], votes - 1}
	if s != wantConflict {
		t.Errorf("s == %v, want %v", s, wantConflict)
	}

	s = k.Vote(groups[0], users[1], targets[0], timeouts[0])
	wantVote := Status{Vote, targets[0], votes - 2}
	if s != wantVote {
		t.Errorf("s == %v, want %v", s, wantVote)
	}

	s = k.Vote(groups[0], users[2], targets[0], timeouts[0])
	wantSuccess := Status{Success, targets[0], 0}
	if s != wantSuccess {
		t.Errorf("s == %v, want %v", s, wantSuccess)
	}

	time.Sleep(duration + delta)
	if flags[0] {
		t.Errorf("flags[0] == %v, want %v", flags[0], false)
	}
	if !flags[1] {
		t.Errorf("flags[1] == %v, want %v", flags[1], true)
	}
}
