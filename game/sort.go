package game

import "sort"

type UserSortFunc func(a, b *User) bool

func (w *World) SortedUsers(f UserSortFunc) []*User {
	users := make([]*User, len(w.Users))
	copy(users, w.Users)
	sort.Slice(users, func(i, j int) bool {
		return f(users[i], users[j])
	})
	return users
}

func ByStrength(a, b *User) bool {
	return a.Strength() > b.Strength()
}

func ByElo(a, b *User) bool {
	return a.Rating > b.Rating
}

func ByWealth(a, b *User) bool {
	return a.Total() > b.Total()
}
