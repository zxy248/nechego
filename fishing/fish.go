package fishing

import "fmt"

type Fish struct {
	Species
	Weight float64 // kilograms
	Length float64 // meters
}

func RandomFish() Fish {
	s := RandomSpecies()
	w := s.randomWeight()
	l := s.Constitution().randomLength(w)
	return Fish{
		Species: s,
		Weight:  w,
		Length:  l,
	}
}

func (f Fish) Price() int {
	return int(f.Weight * f.PricePerKg())
}

func (f Fish) Light() bool {
	return f.Weight < f.NormalWeight()
}

func (f Fish) Heavy() bool {
	return f.Weight > f.NormalWeight()+f.weightSigma()
}

func (f Fish) Cheap() bool {
	return f.Price() < CheapThreshold
}

func (f Fish) Expensive() bool {
	return f.Price() > ExpensiveThreshold
}

func (f Fish) String() string {
	var length, weight string
	if f.Length < 1.0 {
		length = fmt.Sprintf("%.1f см", f.Length*1e2)
	} else {
		length = fmt.Sprintf("%.2f м", f.Length)
	}
	if f.Weight < 1.0 {
		weight = fmt.Sprintf("%.1f г", f.Weight*1e3)
	} else {
		weight = fmt.Sprintf("%.2f кг", f.Weight)
	}
	return fmt.Sprintf("🐟 %s (%s, %s)", f.Species, weight, length)
}

type Fishes []Fish

func (f Fishes) Price() int {
	sum := 0
	for _, ff := range f {
		sum += ff.Price()
	}
	return sum
}

func (f Fishes) Weight() float64 {
	sum := 0.0
	for _, ff := range f {
		sum += ff.Weight
	}
	return sum
}
