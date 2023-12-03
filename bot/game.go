package main

import (
	"math/rand"
	"nechego/game"
	"nechego/item"
	"nechego/money"
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

func onBuyHandler(w *game.World) func(*game.User, *game.Product) {
	return func(u *game.User, p *game.Product) {
		payMarketEmployee(w, p.Price/5)
	}
}

func onSellHandler(w *game.World) func(*game.User, *item.Item, int) {
	return func(u *game.User, i *item.Item, profit int) {
		payTopElo(w, profit/10)
	}
}

func payTopElo(w *game.World, n int) {
	if n == 0 {
		return
	}
	u := w.TopUser(game.ByElo)
	x := item.New(&money.Cash{Money: n})
	u.Funds.Add("налог", x)
}

func payMarketEmployee(w *game.World, n int) {
	if n == 0 {
		return
	}
	if id, ok := w.Market.Shift.Employee(); ok {
		x := item.New(&money.Cash{Money: n})
		w.User(id).Funds.Add("магазин", x)
	}
}
