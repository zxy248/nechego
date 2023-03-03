package game

import "sort"

// Friends is a set of users' IDs.
type Friends map[int64]bool

// Add adds the specified user to Friends.
func (f Friends) Add(x IDer) {
	f[x.ID()] = true
}

// Remove removes the specified user from Friends.
// Returns true if the given user was a friend, or false if not.
func (f Friends) Remove(x IDer) bool {
	defer delete(f, x.ID())
	return f.With(x)
}

// With returns true if the specified user is in Friends.
func (f Friends) With(x IDer) bool {
	return f[x.ID()]
}

// List returns a sorted list of IDs.
func (f Friends) List() []int64 {
	l := []int64{}
	for id := range f {
		l = append(l, id)
	}
	sort.Slice(l, func(i, j int) bool {
		return l[i] < l[j]
	})
	return l
}

// MutualFriends returns true if both users are friends to each other.
func MutualFriends(u, v *User) bool {
	return u.Friends.With(v) && v.Friends.With(u)
}
