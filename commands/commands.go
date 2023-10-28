package commands

import (
	"math/rand"
	"strings"
)

type Command struct {
	Message string
	Photo   string
}

func (c Command) HasPhoto() bool {
	return c.Photo != ""
}

type Commands map[string][]Command

func (cs Commands) Add(s string, c Command) {
	cs[s] = append(cs[s], c)
}

func (cs Commands) Remove(s string) {
	delete(cs, s)
}

func (cs Commands) Match(s string) (c Command, ok bool) {
	var n int
	var r []Command
	for def, cmds := range cs {
		if len(def) > n && strings.Contains(s, def) {
			r = cmds
			n = len(def)
		}
	}
	if r != nil {
		return r[rand.Intn(len(r))], true
	}
	return Command{}, false
}
