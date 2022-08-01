package service

import (
	"sync"
	"time"
)

func PeriodicallyRun(f func(), p time.Duration) func() time.Time {
	next := &muTime{time.Now().Add(p), &sync.RWMutex{}}
	go func() {
		for t := range time.Tick(p) {
			next.set(t.Add(p))
			f()
		}
	}()
	return func() time.Time {
		return next.get()
	}
}

type muTime struct {
	t time.Time
	*sync.RWMutex
}

func (e *muTime) set(t time.Time) {
	e.Lock()
	defer e.Unlock()
	e.t = t
}

func (e *muTime) get() time.Time {
	e.RLock()
	defer e.RUnlock()
	return e.t
}
