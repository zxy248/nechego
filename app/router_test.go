package app

import (
	"nechego/input"
	"testing"
)

func TestCommandHandler(t *testing.T) {
	a := newTestApp()
	for _, c := range input.AllCommands() {
		if c == input.CommandUnknown {
			if h := a.commandHandler(c); h != nil {
				t.Error("want nil")
			}
		} else {
			if h := a.commandHandler(c); h == nil {
				t.Error("want not nil")
			}
		}
	}
}
