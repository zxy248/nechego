package fishing

import "testing"

func TestRod(t *testing.T) {
	const sample = 100_000
	levels := map[string]bool{}
	for i := 0; i < sample; i++ {
		rod := NewRod()
		d := rod.Durability
		if d < 0.5 || d > 1.0 {
			t.Errorf("dur == %v, want in [0.5, 1.0]", d)
		}
		q := rod.Quality
		if q < 0 || q > 1 {
			t.Errorf("q == %v, want in [0, 1]", q)
		}
		levels[rod.level()] = true
	}
	const want = 10
	if len(levels) != want {
		t.Errorf("len(levels) == %v, want %v", len(levels), want)
	}
}
