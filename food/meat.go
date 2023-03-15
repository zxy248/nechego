package food

import (
	"fmt"
	"nechego/pets"
)

type Meat struct {
	Species pets.Species
}

func (m Meat) String() string {
	return fmt.Sprintf("🥩 Мясо (%s)", m.Species)
}
