package fishing

import "testing"

func TestSpeciesData(t *testing.T) {
	i := 0
	for _, s := range species {
		i++
		if s.name == "" {
			t.Errorf("empty string")
		}
		if s.normalWeight <= 0 {
			t.Errorf("normal weight must be positive")
		}
		if s.maximumWeight <= 0 {
			t.Errorf("maximum weight must be positive")
		}
		if s.normalWeight >= s.maximumWeight {
			t.Errorf("maximum weight must be greater than normal weight")
		}
		if s.pricePerKg <= 0 {
			t.Errorf("price must be positive")
		}
	}
	if i != int(numberOfSpecies) {
		t.Errorf("number of elements in speciesData must be equal to %v", numberOfSpecies)
	}
}
