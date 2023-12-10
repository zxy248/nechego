package fishing

import "testing"

func TestNewRod(t *testing.T) {
	const sample = 100_000
	for i := 0; i < sample; i++ {
		r := NewRod()
		if r.Durability < 0.5 || r.Durability > 1.0 {
			t.Errorf("r.Durability = %v, want in [0.5, 1.0]", r.Durability)
		}
	}
}
