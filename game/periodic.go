package game

import (
	"math/rand"
	"time"
)

// startPeriodicTasks runs a set of goroutines corresponding to each
// periodic task.
func (w *World) startPeriodicTasks() {
	table := []struct {
		period time.Duration
		action func(*World)
	}{
		{time.Minute, refillMarket},
		{time.Minute, restoreEnergy},
		{time.Minute, fillNets},
	}
	for _, v := range table {
		go periodically(w, v.period, v.action)
	}
}

// periodically runs the given function at the specified period of time.
func periodically(w *World, p time.Duration, f func(*World)) {
	for range time.NewTicker(p).C {
		w.Lock()
		f(w)
		w.Unlock()
	}
}

// refillMarket adds a new random product to the market.
func refillMarket(w *World) {
	w.Market.Refill()
}

// restoreEnergy adds some energy to all users.
func restoreEnergy(w *World) {
	const delta = 0.01
	for _, u := range w.Users {
		if u.InventoryFull() {
			u.Energy.Add(delta / 2)
		} else {
			u.Energy.Add(delta)
		}
	}
}

// fillNets fills the cast fishing nets.
func fillNets(w *World) {
	for _, u := range w.Users {
		if rand.Float64() < 0.04 {
			u.FillNet()
		}
	}
}
