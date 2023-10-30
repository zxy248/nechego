package services

import (
	"fmt"
	"nechego/chat"
	"nechego/game"
)

type Handler interface {
	Match(*chat.Message) Request
}

type Request interface {
	Process() error
}

func GetWorld(u *game.Universe, id int64) *game.World {
	w, err := u.World(id)
	if err != nil {
		panic(fmt.Sprintf("cannot get world %v", id))
	}
	return w
}

func Do(w *game.World, f func(*game.World)) {
	w.Lock()
	defer w.Unlock()
	f(w)
}
