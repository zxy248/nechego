package game

import "time"

// Shift represents a work shift.
type Shift struct {
	LastEmployee int64
	From, To     time.Time
}

// Employee returns the current employee's ID iff the shift is
// in progress.
func (s *Shift) Employee() (id int64, ok bool) {
	if !s.active() {
		return 0, false
	}
	return s.LastEmployee, true
}

// Start starts the shift and returns true.
// If another shift is in progress, returns false.
func (s *Shift) Start(employee int64, d time.Duration) bool {
	if s.active() {
		return false
	}
	s.LastEmployee = employee
	t := time.Now()
	s.From = t
	s.To = t.Add(d)
	return true
}

func (s *Shift) active() bool {
	return time.Now().Before(s.To)
}
