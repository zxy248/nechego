package game

import "time"

// Shift represents a work shift.
type Shift struct {
	LastWorkerID int64
	Start        time.Time
	End          time.Time
}

// Worker returns the current worker's ID if the shift is in progress.
func (s *Shift) Worker() (id int64, ok bool) {
	if time.Now().After(s.End) {
		return 0, false
	}
	return s.LastWorkerID, true
}

// Begin starts the shift and returns true.
// If another shift is in progress, returns false.
func (s *Shift) Begin(workerID int64, d time.Duration) bool {
	if time.Now().Before(s.End) {
		return false
	}
	s.LastWorkerID = workerID
	s.Start = time.Now()
	s.End = s.Start.Add(d)
	return true
}

// Cancel stops the shift by setting its end to the current time.
func (s *Shift) Cancel() {
	if time.Now().After(s.End) {
		return
	}
	s.End = time.Now()
}
