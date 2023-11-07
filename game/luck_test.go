package game

import (
	"math"
	"nechego/dates"
	"testing"
)

func TestLuck(t *testing.T) {
	const sample = 10000
	lucks := make([]float64, sample)
	sum := 0.0
	min := 1.0
	max := 0.0
	for i := 0; i < sample; i++ {
		u := &User{ID: int64(i)}
		l := u.Luck()
		if l < 0 || l >= 1 {
			t.Errorf("l = %v, want [0, 1)", l)
		}
		if l < min {
			min = l
		}
		if l > max {
			max = l
		}
		lucks[i] = l
		sum += l
	}
	avg := sum / sample

	t.Run("distribution", func(t *testing.T) {
		const (
			n       = 10
			epsilon = 200
			want    = sample / n
		)
		buckets := make([][]float64, n)
		for _, l := range lucks {
			i := int(l * n)
			buckets[i] = append(buckets[i], l)
		}
		for _, b := range buckets {
			if math.Abs(float64(len(b))-sample/n) > epsilon {
				t.Errorf("len(b) = %v, want %v±%v", len(b), want, epsilon)
			}
		}
	})
	t.Run("average", func(t *testing.T) {
		const (
			want    = 0.5
			epsilon = 0.01
		)
		if math.Abs(avg-want) > epsilon {
			t.Errorf("avg = %v, want %v±%v", avg, want, epsilon)
		}
	})
	t.Run("distance", func(t *testing.T) {
		const want = 0.98
		diff := max - min
		if diff < want {
			t.Errorf("diff = %v, want >= %v", diff, want)
		}
	})
	t.Run("delta", func(t *testing.T) {
		const (
			n       = 30
			id      = 109692644
			epsilon = 1e-4
		)
		lucks := make([]float64, n)
		for i := range lucks {
			date := dates.Today().AddDate(0, 0, i)
			lucks[i] = luck(date, id)
		}
		for i, x := range lucks {
			for j, y := range lucks {
				if i == j {
					continue
				}
				delta := math.Abs(x - y)
				if math.Abs(delta) < epsilon {
					t.Errorf("delta = %v, want < %v", delta, epsilon)
				}
			}
		}
	})
}
