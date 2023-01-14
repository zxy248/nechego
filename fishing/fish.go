package fishing

import "fmt"

type Fish struct {
	Species
	Weight float64 // kilograms
	Length float64 // meters
}

func RandomFish() *Fish {
	s := RandomSpecies()
	w := s.randomWeight()
	l := s.Constitution().randomLength(w)
	return &Fish{s, w, l}
}

func (f *Fish) Price() float64  { return f.Weight * f.PricePerKg() }
func (f *Fish) Light() bool     { return f.Weight < f.NormalWeight() }
func (f *Fish) Heavy() bool     { return f.Weight > f.NormalWeight()+f.weightStdDev() }
func (f *Fish) Cheap() bool     { return f.Price() < cheapThreshold }
func (f *Fish) Expensive() bool { return f.Price() > expensiveThreshold }
func (f *Fish) String() string {
	var length, weight string
	if f.Length < 1.0 {
		length = fmt.Sprintf("%.1f см", f.Length*100)
	} else {
		length = fmt.Sprintf("%.2f м", f.Length)
	}
	if f.Weight < 1.0 {
		weight = fmt.Sprintf("%.1f г", f.Weight*1000)
	} else {
		weight = fmt.Sprintf("%.2f кг", f.Weight)
	}
	return fmt.Sprintf("🐟 %s (%s, %s)", f.Species, weight, length)
}
