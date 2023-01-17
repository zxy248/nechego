package game

import (
	"math"
	"testing"
)

func TestLuck(t *testing.T) {
	const sample = 10000
	lucks := make([]float64, sample)
	sum := 0.0
	for i := 0; i < sample; i++ {
		u := &User{TUID: int64(i)}
		l := u.Luck()
		if l < 0 || l >= 1 {
			t.Errorf("l == %v, want [0, 1)", l)
		}
		lucks[i] = l
		sum += l
	}
	avg := sum / sample

	t.Run("distribution", func(t *testing.T) {
		const n, epsilon = 10, 200
		const want = sample / n
		buckets := make([][]float64, n)
		for _, l := range lucks {
			i := int(l * n)
			buckets[i] = append(buckets[i], l)
		}
		for _, b := range buckets {
			if math.Abs(float64(len(b))-sample/n) > epsilon {
				t.Errorf("len(b) == %v, want %v±%v", len(b), want, epsilon)
			}
		}
	})
	t.Run("average", func(t *testing.T) {
		const want, epsilon = 0.5, 0.01
		if math.Abs(avg-want) > epsilon {
			t.Errorf("avg == %v, want %v±%v", avg, want, epsilon)
		}
	})
}
