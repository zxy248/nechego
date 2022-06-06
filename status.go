package main

import "sync"

// status represents the current state of the bot. on is true if the bot is on,
// and false if it is off.
type status struct {
	mu *sync.Mutex
	on bool
}

// newStatus returns a new status. The bot is turned on by default.
func newStatus() *status {
	return &status{mu: &sync.Mutex{}, on: true}
}

// active returns the current state of the bot.
func (s *status) active() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.on
}

// turnOn turns the bot on.
func (s *status) turnOn() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.on = true
}

// turnOff turns the bot off.
func (s *status) turnOff() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.on = false
}
