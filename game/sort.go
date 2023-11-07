package game

import "sort"

type UserSortFunc func(u1, u2 *User) bool

func (w *World) ListUsers() []*User {
	us := make([]*User, 0, len(w.Users))
	for _, u := range w.Users {
		us = append(us, u)
	}
	return us
}

func (w *World) SortedUsers(f UserSortFunc) []*User {
	us := w.ListUsers()
	sort.Slice(us, func(i, j int) bool {
		return f(us[i], us[j])
	})
	return us
}

func (w *World) Position(u *User, f UserSortFunc) int {
	us := w.SortedUsers(f)
	for i, v := range us {
		if u == v {
			return i
		}
	}
	panic("cannot determine user position in sorted list")
}

func ByStrength(u1, u2 *User) bool {
	return u1.Strength() > u2.Strength()
}

func ByElo(u1, u2 *User) bool {
	return u1.Rating > u2.Rating
}

func ByWealth(u1, u2 *User) bool {
	return u1.Balance().Total() > u2.Balance().Total()
}
