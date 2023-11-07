package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"nechego/commands"
	"nechego/fishing"
	"nechego/item"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type World struct {
	ID       int64
	Users    map[int64]*User
	Floor    *item.Set
	Market   *Market
	Casino   *Casino
	Messages int
	History  *fishing.History
	Commands commands.Commands
	Inactive bool

	sync.Mutex `json:"-"`
}

func NewWorld(id int64) *World {
	return &World{
		ID:       id,
		Users:    map[int64]*User{},
		Floor:    item.NewSet(),
		Market:   NewMarket(),
		Casino:   &Casino{Timeout: time.Second * 25},
		History:  fishing.NewHistory(),
		Commands: commands.Commands{},
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

func (w *World) RandomUserID() int64 {
	ids := w.uids()
	return ids[rand.Intn(len(ids))]
}

func (w *World) RandomUserIDs(n int) []int64 {
	ids := w.uids()
	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})
	return ids[:min(len(ids), n)]
}

func (w *World) uids() []int64 {
	ids := make([]int64, 0, len(w.Users))
	for id := range w.Users {
		ids = append(ids, id)
	}
	return ids
}

func (w *World) User(id int64) *User {
	u := w.userByID(id)
	u.ReputationFactor = w.reputationFactor(u.Reputation.Score())
	u.RatingPosition = w.Position(u, ByElo)
	u.Activity = float64(u.Messages) / float64(w.Messages)
	return u
}

func (w *World) userByID(id int64) *User {
	for _, u := range w.Users {
		if u.ID == id {
			return u
		}
	}
	u := NewUser(id)
	w.Users[u.ID] = u
	return u
}

func (w *World) Capital() (total, avg int) {
	for _, u := range w.Users {
		total += u.Balance().Total()
	}
	return total, total / len(w.Users)
}

func (w *World) reputationFactor(n int) float64 {
	low, high := w.reputationBounds()
	d := high - low
	if d == 0 {
		return 0.5
	}
	return float64(n-low) / float64(d)
}

func (w *World) reputationBounds() (low, high int) {
	for _, u := range w.Users {
		s := u.Reputation.Score()
		low, high = min(low, s), max(high, s)
	}
	return
}
