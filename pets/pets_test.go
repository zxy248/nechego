package pets

import (
	"testing"
)

func TestRandom(t *testing.T) {
	const sample = 5000
	males, females := 0, 0
	animals := map[Species]bool{}
	for i := 0; i < sample; i++ {
		pet := Random()
		if pet.Gender == Male {
			males++
		}
		if pet.Gender == Female {
			females++
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
		const want = 100
		delta := males - females
		if delta < 0 {
			delta = -delta
		}
		if delta > want {
			t.Errorf("delta == %d, want <= %d", delta, want)
		}
	})
	t.Run("binary", func(t *testing.T) {
		sum := males + females
		if sum != sample {
			t.Errorf("sum == %d, want %d", sum, sample)
		}
	})
}

func TestData(t *testing.T) {
	for _, v := range species {
		p := v.Probability
		if p <= 0 || p > 1 {
			t.Errorf("p == %v, want in (0, 1]", p)
		}
	}
}
