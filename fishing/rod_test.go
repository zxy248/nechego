package fishing

import (
	"math/rand"
	"testing"
)

func TestNewRod(t *testing.T) {
	const sample = 100_000
	levels := map[string]bool{}
	for i := 0; i < sample; i++ {
		r := NewRod()
		if r.Durability < 0.5 || r.Durability > 1.0 {
			t.Errorf("r.Durability = %v, want in [0.5, 1.0]", r.Durability)
		}
		if r.Quality < 0 || r.Quality > 1 {
			t.Errorf("r.Quality = %v, want in [0, 1]", r.Quality)
		}
		levels[r.level()] = true
	}
	const want = 10
	if len(levels) != want {
		t.Errorf("len(levels) = %v, want %v", len(levels), want)
	}
}

func TestGreek(t *testing.T) {
	const sample = 100_000
	levels := map[string]bool{}
	for i := 0; i < sample; i++ {
		r := NewRod()
		r.Quality = 1 + 2.4*rand.Float64()
		levels[r.level()] = true
	}
	const want = 24
	if len(levels) != want {
		t.Errorf("len(levels) = %v, want %v", len(levels), want)
	}
	for _, g := range greeks {
		if !levels[g] {
			t.Errorf("%v not in levels", g)
		}
	}
}
