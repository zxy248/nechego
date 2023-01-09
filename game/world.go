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
	TGID         int64
	Users        []*User
	Floor        []*Item
	floorHotkeys map[int]*Item
	Market       *Market
	NextItemID   int

	mu sync.Mutex
}

func (w *World) Lock() {
	w.mu.Lock()
}

func (w *World) Unlock() {
	w.mu.Unlock()
}

func NewWorld(id int64) *World {
	return &World{
		TGID:  id,
		Users: []*User{},
		Floor: []*Item{},
	}
}

func LoadWorld(path string) (*World, error) {
	f, err := os.Open(path)
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

func (w *World) Save(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	w.mu.Lock()
	defer w.mu.Unlock()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	return enc.Encode(w)
}

func (w *World) Identify(i *Item) {
	if i.Type == ItemTypeUnknown {
		panic("cannot itemize: item type unknown")
	}
	if i.ID != 0 {
		return
	}
	i.ID = w.NextItemID
	w.NextItemID++
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
	w.TraverseUsers(func(v *User) {
		if v.TUID == tuid {
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

func (w *World) AddItem(u *User, i *Item) {
	w.Identify(i)
	u.Inventory = append(u.Inventory, i)
}

func (w *World) Drop(u *User, i *Item) bool {
	if !i.Transferable {
		return false
	}
	for n, j := range u.Inventory {
		if i == j {
			u.Inventory[n] = u.Inventory[len(u.Inventory)-1]
			u.Inventory = u.Inventory[:len(u.Inventory)-1]
			w.Floor = append(w.Floor, j)
			return true
		}
	}
	return false
}

func (w *World) Pick(u *User, i *Item) bool {
	for n, j := range w.Floor {
		if i == j {
			w.Floor[n] = w.Floor[len(w.Floor)-1]
			w.Floor = w.Floor[:len(w.Floor)-1]
			u.Inventory = append(w.Floor, j)
			return true
		}
	}
	return false
}

func (w *World) ListFloor() []*Item {
	var r []*Item
	w.floorHotkeys, r = hotkeys(w.Floor)
	return r
}

func (w *World) OnFloor(i *Item) bool {
	for _, j := range w.Floor {
		if i == j {
			return true
		}
	}
	return false
}

func (w *World) ItemByKey(k int) (i *Item, ok bool) {
	i, ok = w.floorHotkeys[k]
	if !ok || !w.OnFloor(i) {
		return nil, false
	}
	return i, true
}
