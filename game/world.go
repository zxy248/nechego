package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Universe struct {
	worlds map[int64]*World
	dir    string
	mu     sync.Mutex
}

func NewUniverse(dir string) *Universe {
	return &Universe{
		dir:    dir,
		worlds: map[int64]*World{},
	}
}

func (u *Universe) worldPath(id int64) string {
	return filepath.Join(u.dir, fmt.Sprintf("world%d.json", id))
}

func (u *Universe) MustWorld(id int64) *World {
	w, err := u.World(id)
	if err != nil {
		panic(err)
	}
	return w
}

func (u *Universe) ForEachWorld(action func(*World)) {
	for _, w := range u.worlds {
		w.Lock()
		action(w)
		w.Unlock()
	}
}

func (u *Universe) World(id int64) (*World, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	w, ok := u.worlds[id]
	if !ok {
		p := u.worldPath(id)
		if err := os.MkdirAll(filepath.Dir(p), 0777); err != nil {
			return nil, err
		}
		w, err := LoadWorld(p)
		if errors.Is(err, os.ErrNotExist) {
			w = NewWorld(id)
		} else if err != nil {
			return nil, err
		}
		u.worlds[id] = w
		return w, nil
	}
	return w, nil
}

func (u *Universe) SaveAll() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	for _, w := range u.worlds {
		if err := w.Save(u.worldPath(w.TGID)); err != nil {
			return err
		}
	}
	return nil
}

type World struct {
	TGID     int64
	Users    []*User
	Floor    *Items
	Market   *Market
	Messages int

	sync.Mutex `json:"-"`
}

func NewWorld(id int64) *World {
	return &World{
		TGID:   id,
		Users:  []*User{},
		Floor:  NewItems(),
		Market: NewMarket(),
	}
}

func LoadWorld(name string) (*World, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	w := &World{}
	if err := json.NewDecoder(f).Decode(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *World) Save(name string) error {
	os.Rename(name, fmt.Sprintf("%s-%d", name, time.Now().Unix()))
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	w.Lock()
	defer w.Unlock()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	return enc.Encode(w)
}

func (w *World) AddUser(u *User) {
	w.Users = append(w.Users, u)
}

func (w *World) RandomUser() *User {
	return w.Users[rand.Intn(len(w.Users))]
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

func (w *World) UserByID(tuid int64) (u *User, ok bool) {
	for _, u = range w.Users {
		if u.TUID == tuid {
			return u, true
		}
	}
	return nil, false
}

func (w *World) RestoreEnergy() {
	for _, u := range w.Users {
		u.RestoreEnergy(1)
	}
}

func (w *World) Capital() (total, avg int) {
	for _, w := range w.Users {
		total += w.Total()
	}
	return total, total / len(w.Users)
}
