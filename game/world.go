package game

import (
	"encoding/json"
	"fmt"
	"github.com/zxy248/nechego/commands"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

type World struct {
	ID       int64
	Users    []int64
	Daily    Daily
	Commands commands.Commands
	Inactive bool

	sync.Mutex `json:"-"`
}

type Daily struct {
	Eblan   int64
	Admin   int64
	Pair    [2]int64
	Updated time.Time
}

func NewWorld(id int64) *World {
	return &World{
		ID:       id,
		Users:    []int64{},
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

func (w *World) RandomUser() int64 {
	return w.Users[rand.Intn(len(w.Users))]
}

func (w *World) RandomUsers(n int) []int64 {
	ids := slices.Clone(w.Users)
	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})
	return ids[:min(len(ids), n)]
}
