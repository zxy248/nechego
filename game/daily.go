package game

import (
	"nechego/item"
	"nechego/token"
)

func (w *World) DailyEblan() (u *User, ok bool) {
	for _, u = range w.Users {
		if u.Eblan() {
			return u, true
		}
	}
	return w.rollDailyEblan(), true
}

func (w *World) rollDailyEblan() *User {
	u := w.User(w.RandomUserID())
	u.Inventory.Add(item.New(&token.Eblan{}))
	return u
}

func (w *World) DailyAdmin() (u *User, ok bool) {
	for _, u = range w.Users {
		if u.Admin() {
			return u, true
		}
	}
	return w.rollDailyAdmin(), true
}

func (w *World) rollDailyAdmin() *User {
	u := w.User(w.RandomUserID())
	u.Inventory.Add(item.New(&token.Admin{}))
	return u
}

func (w *World) DailyPair() (pair []*User, ok bool) {
	if len(w.Users) < 2 {
		return nil, false
	}
	for _, u := range w.Users {
		if u.Pair() {
			pair = append(pair, u)
		}
		if len(pair) == 2 {
			break
		}
	}
	if len(pair) != 2 {
		return w.rollDailyPair()
	}
	return pair, true
}

func (w *World) rollDailyPair() (pair []*User, ok bool) {
	p := w.RandomUserIDs(2)
	if len(p) != 2 {
		return nil, false
	}
	u1 := w.User(p[0])
	u2 := w.User(p[1])
	u1.Inventory.Add(item.New(&token.Pair{}))
	u2.Inventory.Add(item.New(&token.Pair{}))
	return []*User{u1, u2}, true
}
