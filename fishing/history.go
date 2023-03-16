package fishing

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

// RecordAnnouncer is a function that is called when a new fishing
// record is set.
type RecordAnnouncer func(e *Entry, p Parameter)

// Parameter of a fish.
type Parameter int

const (
	Weight Parameter = iota
	Length
	Price
)

var parameters = [...]Parameter{Weight, Length, Price}

// Entry of a history.
type Entry struct {
	TUID int64
	Fish *Fish
	Time time.Time
}

// History of caught fish.
type History struct {
	entries  []*Entry
	top      map[Parameter]*Entry
	announce RecordAnnouncer
	mu       sync.Mutex
}

// NewHistory returns a new History.
func NewHistory() *History {
	return &History{entries: []*Entry{}}
}

// MarshalJSON implements the json.Marshaler interface.
func (h *History) MarshalJSON() ([]byte, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	return json.Marshal(h.entries)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (h *History) UnmarshalJSON(data []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := json.Unmarshal(data, &h.entries); err != nil {
		return err
	}
	h.rebuildTop()
	return nil
}

// rebuildTop rebuilds the history top from the list of entries.
func (h *History) rebuildTop() {
	for _, e := range h.entries {
		for _, p := range parameters {
			h.contend(e, p)
		}
	}
}

// Add adds a new entry to the history.
func (h *History) Add(tuid int64, f *Fish) {
	go func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		e := &Entry{tuid, f, time.Now()}
		h.entries = append(h.entries, e)
		for _, p := range parameters {
			h.contend(e, p)
		}
	}()
}

// Top returns the first n records of the parameter p.
func (h *History) Top(p Parameter, n int) []*Entry {
	h.mu.Lock()
	defer h.mu.Unlock()

	r := make([]*Entry, len(h.entries))
	copy(r, h.entries)
	sort.Slice(r, func(i, j int) bool {
		switch p {
		case Weight:
			return r[i].Fish.Weight > r[j].Fish.Weight
		case Length:
			return r[i].Fish.Length > r[j].Fish.Length
		case Price:
			return r[i].Fish.Price() > r[j].Fish.Price()
		}
		panic(fmt.Sprintf("unexpected parameter %v", p))
	})
	if len(r) < n {
		return r
	}
	return r[:n]
}

// Announce sets the given function to be called on new records.
func (h *History) Announce(f RecordAnnouncer) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.announce = f
}

// contend compares the history entry e with the top record of the
// parameter p. Announces a new record if the given entry makes it to
// the top and the record announcer is set.
func (h *History) contend(e *Entry, p Parameter) {
	if h.top == nil {
		h.top = map[Parameter]*Entry{}
	}
	if t, ok := h.top[p]; !ok || param(e.Fish, p) > param(t.Fish, p) {
		h.top[p] = e
		if h.announce != nil {
			go h.announce(e, p)
		}
	}
}

// param returns the parameter p of the fish f.
func param(f *Fish, p Parameter) float64 {
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
