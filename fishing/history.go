package fishing

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

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
	entries []*Entry
	top     map[Parameter]*Entry
	records map[Parameter]chan *Entry
	mu      sync.Mutex
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
	for _, e := range h.entries {
		for _, p := range parameters {
			h.contend(e, p)
		}
	}
	return nil
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

	entires := make([]*Entry, len(h.entries))
	copy(entires, h.entries)
	sort.Slice(entires, func(i, j int) bool {
		switch p {
		case Weight:
			return entires[i].Fish.Weight > entires[j].Fish.Weight
		case Length:
			return entires[i].Fish.Length > entires[j].Fish.Length
		case Price:
			return entires[i].Fish.Price() > entires[j].Fish.Price()
		}
		panic(fmt.Sprintf("unexpected parameter %v", p))
	})
	if len(entires) > n {
		return entires[:n]
	}
	return entires
}

// Records returns a new channel of the given parameter for record
// announcements. Panics if the channel is already created.
func (h *History) Records(p Parameter) chan *Entry {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.records == nil {
		h.records = make(map[Parameter]chan *Entry)
	}
	if h.records[p] != nil {
		panic(fmt.Sprintf("record channel %v is already created", p))
	}
	h.records[p] = make(chan *Entry)
	return h.records[p]
}

// contend compares the history entry e with the top record of the
// given parameter p. If it is a new record, sends a corresponding
// value on the record channel of p.
func (h *History) contend(e *Entry, p Parameter) {
	if h.top == nil {
		h.top = make(map[Parameter]*Entry)
	}
	t, ok := h.top[p]
	if !ok || param(e.Fish, p) > param(t.Fish, p) {
		h.top[p] = e
		h.announce(e, p)
	}
}

// announce sends the entry e on the record channel of the given parameter.
func (h *History) announce(e *Entry, p Parameter) {
	select {
	case h.records[p] <- e:
	default:
		// Channel is nil; drop the value.
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
