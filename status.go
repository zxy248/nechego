package main

import "sync"

type status struct {
	mu *sync.Mutex
	s  bool
}

// newStatus returns a new active status
func newStatus() *status {
	return &status{mu: &sync.Mutex{}, s: true}
}

// active returns the current bot status
func (s *status) active() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.s
}

// turnOn turns the bot on
func (s *status) turnOn() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.s = true
}

// turnOff turn the bot off
func (s *status) turnOff() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.s = false
}
