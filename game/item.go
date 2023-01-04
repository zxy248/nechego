package game

import (
	"fmt"
	"time"
)

const (
	Hotkeys       = "123456789йцукенгшщзхфывапролджэячсмитьбю"
	InventorySize = len(Hotkeys)
)

type Object interface {
	fmt.Stringer
}

type Item struct {
	ID           int
	Transferable bool
	Expire       *time.Time
	Object       Object
}
