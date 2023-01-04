package game

import (
	"encoding/json"
	"math/rand"
	"os"
)

const WorldDir = "world"
const WorldFile = "world.json"

type World struct {
	ID    int
	Users []*User
	Floor []*Item
}

func (w *World) Save() error {
	f, err := os.OpenFile(WorldFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "\t")
	return e.Encode(w)
}

func (w *World) Load() error {
	f, err := os.Open(WorldFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(w)
}

func (w *World) RandomUsers(n int) []*User {
	users := make([]*User, len(w.Users))
	copy(users, w.Users)
	rand.Shuffle(len(users), func(i, j int) {
		users[i], users[j] = users[j], users[i]
	})
	if len(users) < n {
		n = len(users)
	}
	return users[:n]
}

func (w *World) UserByID(id int) (u *User, ok bool) {
	w.TraverseUsers(func(v *User) {
		if v.ID == id {
			u, ok = v, true
			return
		}
	})
	return
}

func (w *World) TraverseUsers(f func(*User)) {
	for _, v := range w.Users {
		f(v)
	}
}

func (w *World) RestoreEnergy() {
	w.TraverseUsers(func(u *User) {
		u.RestoreEnergy(1)
	})
}

func (w *World) DailyEblan() (u *User, ok bool) {
	w.TraverseUsers(func(v *User) {
		if v.IsEblan() {
			u, ok = v, true
			return
		}
	})
	return
}

func (w *World) DailyAdmin() (u *User, ok bool) {
	w.TraverseUsers(func(v *User) {
		if v.IsAdmin() {
			u, ok = v, true
			return
		}
	})
	return
}

func (w *World) DailyPair() (u [2]*User, ok bool) {
	n := 0
	w.TraverseUsers(func(v *User) {
		if n == 2 {
			ok = true
			return
		}
		if v.IsPair() {
			u[n] = v
			n++
		}
	})
	return
}
