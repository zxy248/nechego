package pvp

import (
	"fmt"
	"time"
)

const WaitForPvE = 15 * time.Minute

// Status is either PvE or PvP.
type Status int

const (
	PvE Status = iota
	PvP
)

func (s Status) String() string {
	switch s {
	case PvE:
		return "PvE"
	case PvP:
		return "PvP"
	}
	panic(fmt.Sprintf("unexpected status %d", s))
}

// Combat mode.
// Do not use the fields of this structure directly; they are public
// for serialization purposes.
type Mode struct {
	PvP     bool
	Updated time.Time
}

// Toggle toggles the combat Mode.
func (m *Mode) Toggle() Status {
	m.PvP = !m.PvP
	m.Updated = time.Now()
	if m.PvP {
		return PvP
	}
	return PvE
}

// Status returns the current combat status.
func (m *Mode) Status() Status {
	if !m.PvP && time.Since(m.Updated) > WaitForPvE {
		return PvE
	}
	return PvP
}
