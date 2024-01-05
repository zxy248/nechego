package game

import "time"

func (w *World) updateDaily() {
	t := time.Now()
	if w.Daily.Updated.IsZero() || w.Daily.Updated.Day() != t.Day() {
		w.Daily.Eblan = w.RandomUser()
		w.Daily.Admin = w.RandomUser()
		w.Daily.Pair[0] = w.RandomUser()
		w.Daily.Pair[1] = w.RandomUser()
		w.Daily.Updated = t
	}
}

func (w *World) DailyEblan() int64 {
	w.updateDaily()
	return w.Daily.Eblan
}

func (w *World) DailyAdmin() int64 {
	w.updateDaily()
	return w.Daily.Admin
}

func (w *World) DailyPair() [2]int64 {
	w.updateDaily()
	return w.Daily.Pair
}
