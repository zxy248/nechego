package pets

import "testing"

func TestRandom(t *testing.T) {
	const sample = 10000
	males, females := 0, 0
	rare := 0
	animals := map[Species]bool{}
	for i := 0; i < sample; i++ {
		pet := Random()
		if pet.Gender == Male {
			males++
		}
		if pet.Gender == Female {
			females++
		}
		if pet.Species.Probability() <= 0.1 {
			rare++
		}
		animals[pet.Species] = true
	}
	t.Run("species coverage", func(t *testing.T) {
		if len(animals) != len(species) {
			t.Errorf("got %d/%d, want all species covered",
				len(animals), len(species))
		}
	})
	t.Run("genders", func(t *testing.T) {
		const want = sample * 0.02
		delta := males - females
		if delta < 0 {
			delta = -delta
		}
		if delta > want {
			t.Errorf("delta = %v, want <= %v", delta, want)
		}
	})
	t.Run("binary", func(t *testing.T) {
		sum := males + females
		if sum != sample {
			t.Errorf("sum = %d, want %d", sum, sample)
		}
	})
	t.Run("few rare", func(t *testing.T) {
		const want = 0.1
		q := float64(rare) / sample
		if q > want {
			t.Errorf("q = %v, want <= %v", q, want)
		}
	})
}

func TestData(t *testing.T) {
	for _, v := range species {
		p := v.Probability
		if p <= 0 || p > 1 {
			t.Errorf("p = %v, want in (0, 1]", p)
		}
	}
}
