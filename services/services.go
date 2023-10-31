package services

import (
	"fmt"
	"nechego/chat"
	"nechego/game"
	"regexp"
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

type State struct {
	universe *game.Universe
	groupID  int64
}

func GroupState(u *game.Universe, m *chat.Message) *State {
	return &State{u, m.Group.ID}
}

func (s *State) Do(f func(*game.World)) {
	w := GetWorld(s.universe, s.groupID)
	w.Lock()
	defer w.Unlock()
	f(w)
}

func Regexp(pattern string) *regexp.Regexp {
	return regexp.MustCompile("(?i)" + pattern)
}
