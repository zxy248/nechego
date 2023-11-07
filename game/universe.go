package game

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Universe holds worlds.
type Universe struct {
	worlds map[int64]*World // Loaded worlds indexed by group IDs.
	dir    string           // Persistent storage directory.
	init   func(*World)     // World initialization function.

	mu sync.Mutex
}

// NewUniverse returns a new Universe.
func NewUniverse(dir string, init func(*World)) *Universe {
	return &Universe{
		dir:    dir,
		worlds: map[int64]*World{},
		init:   init,
	}
}

// worldPath returns the location of a save file of the world by the
// specified ID.
func (u *Universe) worldPath(id int64) string {
	return filepath.Join(u.dir, fmt.Sprintf("world%d.json", id))
}

// World returns the world by the given ID from the universe. If the
// world is not active, loads the save from the world's file. If there
// is no save file found, creates a new world.
func (u *Universe) World(id int64) (*World, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	w, ok := u.worlds[id]
	if !ok {
		// Invariant: the world in not initialized.
		// This case holds only once for each world.
		var err error
		w, err = LoadWorld(u.worldPath(id))
		if errors.Is(err, os.ErrNotExist) {
			w = NewWorld(id)
		} else if err != nil {
			return nil, err
		}
		u.init(w)
		u.worlds[id] = w
	}
	return w, nil
}

// SaveAll saves all active worlds.
func (u *Universe) SaveAll() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	for _, w := range u.worlds {
		if err := w.Save(u.worldPath(w.ID)); err != nil {
			return err
		}
	}
	return nil
}
