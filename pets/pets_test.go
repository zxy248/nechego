package pets

import "testing"

func TestRandomSpecies(t *testing.T) {
	values := []float64{0.004, 0.12, 0.23, 0.47, 0.89}
	for _, v := range values {
		var prev, curr Species
		prev = randomSpecies(v)
		for i := 0; i < 10; i++ {
			curr = randomSpecies(v)
			if prev != curr {
				break
			}
			prev = randomSpecies(v)
		}
	}
}
