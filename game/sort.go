package game

import "sort"

type UserSortFunc func(w *World, a, b *User) bool

func (w *World) SortedUsers(f UserSortFunc) []*User {
	users := make([]*User, len(w.Users))
	copy(users, w.Users)
	sort.Slice(users, func(i, j int) bool {
		return f(w, users[i], users[j])
	})
	return users
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
