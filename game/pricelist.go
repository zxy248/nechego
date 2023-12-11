package game

import (
	"math"
	"math/rand"
	"nechego/details"
	"nechego/farm/plant"
	"nechego/fishing"
	"nechego/item"
	"time"
)

// PriceList is a dynamically updated list of prices.
type PriceList struct {
	Updated time.Time
	Plants  map[plant.Type]int
}

// NewPriceList returns a new refreshed price list.
func NewPriceList() *PriceList {
	p := &PriceList{Plants: map[plant.Type]int{}}
	p.Refresh()
	return p
}

// Refresh updates the price list idempotently.
func (p *PriceList) Refresh() {
	if !p.Updated.IsZero() && time.Now().YearDay() == p.Updated.YearDay() {
		return
	}
	p.Updated = time.Now()
	p.refreshPlants()
}

func (p *PriceList) refreshPlants() {
	for _, t := range plant.Types {
		x := t.Price()
		q := math.Abs(x + 0.33*x*rand.NormFloat64())
		p.Plants[t] = int(q)
	}
}

// Price returns the price of the given plant type.
func (p *PriceList) Price(x *item.Item) (price int, ok bool) {
	p.Refresh()
	switch v := x.Value.(type) {
	case *fishing.Fish:
		return p.fishPrice(v), true
	case *details.Details:
		return p.detailsPrice(v), true
	case *plant.Plant:
		return p.plantPrice(v), true
	}
	return 0, false
}

func (p *PriceList) fishPrice(f *fishing.Fish) int {
	return int(f.Price())
}

func (p *PriceList) detailsPrice(d *details.Details) int {
	return 100 * d.Count
}

func (p *PriceList) plantPrice(l *plant.Plant) int {
	return p.Plants[l.Type] * l.Count
}
