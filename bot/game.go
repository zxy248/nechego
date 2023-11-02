package main

import (
	"math/rand"
	"nechego/game"
	"time"
)

func runServices(w *game.World) {
	table := []struct {
		period time.Duration
		action func(*game.World)
	}{
		{time.Minute, refreshMarket},
		{time.Minute, restoreEnergy},
		{time.Minute, fillNets},
	}
	for _, v := range table {
		d := v.period
		f := v.action
		go func() {
			for range time.NewTicker(d).C {
				w.Lock()
				f(w)
				w.Unlock()
			}
		}()
	}
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

func fillNets(w *game.World) {
	for _, u := range w.Users {
		if rand.Float64() < 0.04 {
			u.FillNet()
		}
	}
}
