package kick

import (
	"sync"
	"time"
)

type Event int

const (
	Init Event = iota
	Vote
	Success
	Timeout
	Duplicate
	Conflict
	NoUser
)

type session struct {
	target int64
	voters map[int64]bool
	timer  *time.Timer
}

type Kicker struct {
	sessions        map[int64]*session
	sessionDuration time.Duration
	votesNeeded     int
	mu              sync.Mutex
}

type Status struct {
	Event          Event
	Target         int64
	VotesRemaining int
}

func New(sessionDuration time.Duration, votesNeeded int) *Kicker {
	return &Kicker{
		sessions:        map[int64]*session{},
		sessionDuration: sessionDuration,
		votesNeeded:     votesNeeded,
	}
}

func (k *Kicker) newSession(group, target, voter int64, timeout func()) *session {
	return &session{
		target: target,
		voters: map[int64]bool{voter: true},
		timer: time.AfterFunc(k.sessionDuration, func() {
			timeout()
			k.mu.Lock()
			defer k.mu.Unlock()
			delete(k.sessions, group)
		}),
	}
}

func (k *Kicker) sessionStatus(s *session, e Event) Status {
	return Status{e, s.target, k.votesNeeded - len(s.voters)}
}

func (k *Kicker) Vote(group, voter, target int64, timeout func()) Status {
	k.mu.Lock()
	defer k.mu.Unlock()
	s, ok := k.sessions[group]

	// session not found; do nothing
	if !ok && target == 0 {
		return Status{NoUser, 0, 0}
	}
	// session not found; initialize
	if !ok {
		s = k.newSession(group, target, voter, timeout)
		k.sessions[group] = s
		return k.sessionStatus(s, Init)
	}
	// session is going for another user
	if target != 0 && target != s.target {
		return k.sessionStatus(s, Conflict)
	}
	// voter is already voted
	if s.voters[voter] {
		return k.sessionStatus(s, Duplicate)
	}
	s.voters[voter] = true
	st := k.sessionStatus(s, Vote)
	if st.VotesRemaining == 0 {
		st.Event = Success
		s.timer.Stop()
		delete(k.sessions, group)
	} else {
		s.timer.Reset(k.sessionDuration)
	}
	return st
}
