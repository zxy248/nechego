package reputation

import (
	"time"

	"golang.org/x/exp/slices"
)

type Vote struct {
	ID   int64
	Dir  Dir
	Time time.Time
}

type Reputation []Vote

type Dir int

const (
	Up Dir = iota
	Down
)

func (r Reputation) Update(id int64, d Dir) Reputation {
	return append(slices.Clone(r), Vote{id, d, time.Now()})
}

func (r Reputation) Total() int {
	t := 0
	for _, v := range r {
		switch v.Dir {
		case Up:
			t++
		case Down:
			t--
		}
	}
	return t
}

func (r Reputation) Last(id int64) Vote {
	v := Vote{}
	for _, x := range r {
		if x.ID == id && x.Time.After(v.Time) {
			v = x
		}
	}
	return v
}
