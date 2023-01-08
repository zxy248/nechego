package game

import "time"

func (w *World) DailyEblan() (u *User, ok bool) {
	w.TraverseUsers(func(v *User) {
		if v.IsEblan() {
			u, ok = v, true
			return
		}
	})
	if !ok {
		return w.rollDailyEblan(), true
	}
	return
}

func (w *World) rollDailyEblan() *User {
	u := w.RandomUser()
	w.AddItem(u, &Item{
		Type:   ItemTypeEblanToken,
		Value:  &EblanToken{},
		Expire: tomorrow(),
	})
	return u
}

func (w *World) DailyAdmin() (u *User, ok bool) {
	w.TraverseUsers(func(v *User) {
		if v.IsAdmin() {
			u, ok = v, true
			return
		}
	})
	if !ok {
		return w.rollDailyAdmin(), true
	}
	return
}

func (w *World) rollDailyAdmin() *User {
	u := w.RandomUser()
	w.AddItem(u, &Item{
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
	w.TraverseUsers(func(v *User) {
		if v.IsPair() {
			pair = append(pair, v)
		}
		if len(pair) == 2 {
			return
		}
	})
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
	w.AddItem(pair[0], pairToken())
	w.AddItem(pair[1], pairToken())
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
