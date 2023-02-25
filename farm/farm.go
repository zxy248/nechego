package farm

import (
	"fmt"
	"math/rand"
	"nechego/farm/plant"
	"strings"
	"time"
)

// Plot represents the Crop's position at the Farm.
type Plot struct {
	Row    int
	Column int
}

// MarshalText implements the encoding.TextMarshaler interface.
func (p Plot) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%d, %d", p.Row, p.Column)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (p *Plot) UnmarshalText(text []byte) error {
	var row, column int
	_, err := fmt.Sscanf(string(text), "%d, %d", &row, &column)
	if err != nil {
		return fmt.Errorf("cannot unmarshal Plot: %v", err)
	}
	p.Row = row
	p.Column = column
	return nil
}

// Crop represents the Plant growing at the Farm.
type Crop struct {
	plant.Type
	Grown time.Time
}

// String returns the emoji corresponding to the current grow status.
func (c Crop) String() string {
	switch {
	case c.Ready(), c.Empty():
		return c.Type.String()
	case time.Until(c.Grown) < 2*time.Hour:
		return "🌿"
	default:
		return "🌱"
	}
}

// Empty is true if the Plant is Void.
func (c Crop) Empty() bool {
	return c.Type == plant.Void
}

// Ready is true if the Crop can be harvested.
func (c Crop) Ready() bool {
	return !c.Empty() && time.Now().After(c.Grown)
}

// Farm represents a place where a player can grow vegetables.
type Farm struct {
	Grid    map[Plot]Crop
	Rows    int
	Columns int
}

// New returns an empty farm of the given size.
func New(rows, columns int) *Farm {
	return &Farm{
		Grid:    map[Plot]Crop{},
		Rows:    rows,
		Columns: columns,
	}
}

// String returns the representation of the Farm as a grid.
func (f *Farm) String() string {
	const fence = "🟰"
	var out string
	for r := 0; r < f.Rows; r++ {
		if r > 0 {
			out += "\n"
		}
		var row string
		for c := 0; c < f.Columns; c++ {
			row += f.Grid[Plot{r, c}].String()
		}
		out += fence + row + fence
	}
	wall := strings.Repeat(fence, 2+f.Columns)
	out = wall + "\n" + out + "\n" + wall
	return out
}

// Harvest pops all ready to harvest crops from the Farm and returns
// them as a list of Plants.
func (f *Farm) Harvest() []*plant.Plant {
	harvested := map[plant.Type]int{}
	for r := 0; r < f.Rows; r++ {
		for c := 0; c < f.Columns; c++ {
			p := Plot{r, c}
			if crop := f.Grid[p]; crop.Ready() {
				harvested[crop.Type] += 1 + rand.Intn(5)
				delete(f.Grid, p)
			}
		}
	}
	r := []*plant.Plant{}
	for typ, count := range harvested {
		r = append(r, &plant.Plant{Type: typ, Count: count})
	}
	return r
}

// Pick pops the Plant from the specified location.
func (f *Farm) Pick(q Plot) (p *plant.Plant, ok bool) {
	if crop := f.Grid[q]; crop.Ready() {
		delete(f.Grid, q)
		return &plant.Plant{
			Type:  crop.Type,
			Count: 1 + rand.Intn(5),
		}, true
	}
	return nil, false
}

// Plant adds a new crop to the Farm and returns true if there is
// enough space. If there is not enough space, returns false.
func (f *Farm) Plant(t plant.Type) bool {
	for r := 0; r < f.Rows; r++ {
		for c := 0; c < f.Columns; c++ {
			p := Plot{r, c}
			if crop := f.Grid[p]; crop.Empty() {
				f.Grid[p] = Crop{t, time.Now().Add(4 * time.Hour)}
				return true
			}
		}
	}
	return false
}

// Free returns the number of empty crops at the Farm.
func (f *Farm) Free() int {
	n := 0
	for r := 0; r < f.Rows; r++ {
		for c := 0; c < f.Columns; c++ {
			if f.Grid[Plot{r, c}].Empty() {
				n++
			}
		}
	}
	return n
}

// Size returns the size of the Farm (Rows × Columns).
func (f *Farm) Size() int {
	return f.Rows * f.Columns
}

// Pending returns the number of plants that can be harvested.
func (f *Farm) Pending() int {
	n := 0
	for r := 0; r < f.Rows; r++ {
		for c := 0; c < f.Columns; c++ {
			if f.Grid[Plot{r, c}].Ready() {
				n++
			}
		}
	}
	return n
}

// Grow adds a new row (or column) to the Farm such that the number of
// rows is less then or equal to the number of columns.
func (f *Farm) Grow() {
	if f.Columns > f.Rows {
		f.Rows++
	} else {
		f.Columns++
	}
}
