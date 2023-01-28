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

func (x Meat) Nutrition() float64 {
	switch x.Species.Size() {
	case pets.Small:
		return 0.10
	case pets.Medium:
		return 0.16
	case pets.Big:
		return 0.22
	}
	panic(fmt.Errorf("unexpected meat size %v", x.Species.Size()))
}
