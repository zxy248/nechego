package game

import "sort"

type UserSortFunc func(w *World, a, b *User) bool

func (w *World) SortedUsers(f UserSortFunc) []*User {
	us := make([]*User, len(w.Users))
	copy(us, w.Users)
	sort.Slice(us, func(i, j int) bool {
		return f(w, us[i], us[j])
	})
	return us
}

func ByStrength(w *World, a, b *User) bool {
	return a.Strength(w) > b.Strength(w)
}

func ByElo(w *World, a, b *User) bool {
	return a.Rating > b.Rating
}

func ByWealth(w *World, a, b *User) bool {
	return a.Balance().Total() > b.Balance().Total()
}
