package main

import "sync"

// status represents the current state of the bot.
type status struct {
	mu     *sync.Mutex
	global bool
	local  map[int64]struct{} // the set of groups where the bot is turned off
}

// newStatus returns a new status. The bot is turned on by default.
func newStatus() *status {
	return &status{
		mu:     &sync.Mutex{},
		global: true,
		local:  make(map[int64]struct{}),
	}
}

// activeGlobal returns the global state of the bot.
func (s *status) activeGlobal() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.global
}

// turnOnGlobal turns the bot on globally.
func (s *status) turnOnGlobal() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.global = true
}

// turnOffGlobal turns the bot off globally.
func (s *status) turnOffGlobal() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.global = false
}

// activeLocal returns the local group state of the bot.
func (s *status) activeLocal(groupID int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.local[groupID]
	return !ok
}

// turnOnLocal turns the bot on for the specified group.
func (s *status) turnOnLocal(groupID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.local, groupID)
}

// turnOffLocal turns the bot off for the specified group.
func (s *status) turnOffLocal(groupID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.local[groupID] = struct{}{}
}
