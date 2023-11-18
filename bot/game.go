package main

import (
	"math/rand"
	"nechego/game"
	"time"
)

func addService(w *game.World, f func(*game.World), p time.Duration) {
	go func() {
		for range time.NewTicker(p).C {
			w.Lock()
			f(w)
			w.Unlock()
		}
	}()
}

func refreshMarket(w *game.World) {
	const c = 10
	for i := len(w.Market.Products()); i < c; i++ {
		w.Market.Refill()
	}
	w.Market.Refill()
	w.Market.Trim(c)
}

func restoreEnergy(w *game.World) {
	const d = 0.01
	for _, u := range w.Users {
		u.Energy.Add(d)
	}
}

func resetEnergy(w *game.World) {
	for _, u := range w.Users {
		u.Energy = 1.0
	}
}

func fillNets(w *game.World) {
	const p = 0.04
	for _, u := range w.Users {
		if rand.Float64() < p {
			u.FillNet()
		}
	}
}
