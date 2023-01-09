package game

import "time"

func (w *World) DailyEblan() (u *User, ok bool) {
	for _, u = range w.Users {
		if u.IsEblan() {
			return u, true
		}
	}
	return w.rollDailyEblan(), true
}

func (w *World) rollDailyEblan() *User {
	u := w.RandomUser()
	u.Inventory.Add(&Item{
		Type:   ItemTypeEblanToken,
		Value:  &EblanToken{},
		Expire: tomorrow(),
	})
	return u
}

func (w *World) DailyAdmin() (u *User, ok bool) {
	for _, u = range w.Users {
		if u.IsAdmin() {
			return u, true
		}
	}
	return w.rollDailyAdmin(), true
}

func (w *World) rollDailyAdmin() *User {
	u := w.RandomUser()
	u.Inventory.Add(&Item{
		Type:         ItemTypeAdminToken,
		Value:        &AdminToken{},
		Expire:       tomorrow(),
		Transferable: true,
	})
	return u
}

func (w *World) DailyPair() (pair []*User, ok bool) {
	if len(w.Users) < 2 {
		return nil, false
	}
	for _, u := range w.Users {
		if u.IsPair() {
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
	pair = w.RandomUsers(2)
	if len(pair) != 2 {
		return nil, false
	}
	pair[0].Inventory.Add(pairToken())
	pair[1].Inventory.Add(pairToken())
	return pair, true
}

func pairToken() *Item {
	return &Item{
		Type:   ItemTypePairToken,
		Value:  &PairToken{},
		Expire: tomorrow(),
	}
}

func tomorrow() time.Time {
	return time.Now().Add(time.Hour * 24)
}
