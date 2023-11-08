package fishing

import (
	"fmt"
	"sort"
	"time"
)

// Parameter of a fish.
type Parameter int

const (
	Weight Parameter = iota
	Length
	Price
)

var Parameters = [...]Parameter{Weight, Length, Price}

// Entry of a history.
type Entry struct {
	ID   int64
	Fish *Fish
	Time time.Time
}

// History of caught fish.
type History struct {
	Entries []*Entry
	Records map[Parameter]*Entry
}

func NewHistory() *History {
	return &History{
		Entries: []*Entry{},
		Records: map[Parameter]*Entry{},
	}
}

// Add adds a new entry to the history.
func (h *History) Add(id int64, f *Fish) {
	e := &Entry{id, f, time.Now()}
	h.Entries = append(h.Entries, e)
	for _, p := range Parameters {
		t, ok := h.Records[p]
		if !ok || param(p, e.Fish) > param(p, t.Fish) {
			h.Records[p] = e
		}
	}
}

// Top returns the first n records of the parameter p.
func (h *History) Top(p Parameter, n int) []*Entry {
	r := make([]*Entry, len(h.Entries))
	copy(r, h.Entries)
	sort.Slice(r, func(i, j int) bool {
		return param(p, r[i].Fish) > param(p, r[j].Fish)
	})
	return r[:min(len(r), n)]
}

// param returns the parameter p of the fish f.
func param(p Parameter, f *Fish) float64 {
	switch p {
	case Weight:
		return f.Weight
	case Length:
		return f.Length
	case Price:
		return f.Price()
	}
	panic(fmt.Sprintf("unexpected parameter %v", p))
}
