package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"nechego/auction"
	"nechego/fishing"
	"nechego/item"
	"nechego/phone"
	"nechego/token"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type World struct {
	TGID     int64
	Users    []*User
	Floor    *item.Set
	Auction  *auction.Auction
	Market   *Market
	Casino   *Casino
	Messages int
	SMS      phone.Database
	History  *fishing.History
	Inactive   bool

	sync.Mutex `json:"-"`
}

func NewWorld(id int64) *World {
	return &World{
		TGID:    id,
		Users:   []*User{},
		Floor:   item.NewSet(),
		Auction: auction.New(),
		Market:  NewMarket(),
		Casino:  &Casino{Timeout: time.Second * 25},
		SMS:     phone.Database{},
		History: fishing.NewHistory(),
	}
}

func LoadWorld(name string) (*World, error) {
	if err := os.MkdirAll(filepath.Dir(name), 0777); err != nil {
		return nil, err
	}

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
	const layout = "06-01-02-15-04-05"
	os.Rename(name, fmt.Sprintf("%s-%s", name, time.Now().Format(layout)))
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
	u.Balance().Add(5000)
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
		return users
	}
	return users[:n]
}

func (w *World) UserByID(tuid int64) *User {
	for _, u := range w.Users {
		if u.TUID == tuid {
			return u
		}
	}
	u := NewUser(tuid)
	w.AddUser(u)
	return u
}

func (w *World) Capital() (total, avg int) {
	for _, w := range w.Users {
		total += w.Balance().Total()
	}
	return total, total / len(w.Users)
}

// runMigrations makes the world consistent with new features.
func (w *World) runMigrations() {
	for _, u := range w.Users {
		// Shrink the farm to its maximum size.
		n := 0
		if dr := u.Farm.Rows - MaxFarmRows; dr > 0 {
			u.Farm.Rows = MaxFarmRows
			n += dr
		}
		if dc := u.Farm.Columns - MaxFarmColumns; dc > 0 {
			u.Farm.Columns = MaxFarmColumns
			n += dc
		}
		if n > 0 {
			u.Inventory.Add(item.New(&token.Legacy{Count: n}))
		}

		// Fill nil values.
		if u.Friends == nil {
			u.Friends = Friends{}
		}
	}
}
