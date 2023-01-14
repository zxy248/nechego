package fishing

import "testing"

func TestSpeciesData(t *testing.T) {
	for _, s := range species {
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
}
