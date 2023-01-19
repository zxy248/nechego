package fishing

import (
	"fmt"
	"math/rand"
)

type Rod struct {
	Quality    float64 // [0, 1]
	Durability float64 // [0, 1]
}

func (r Rod) String() string {
	return fmt.Sprintf("🎣 Удочка (%s, %.f%%)", r.level(), r.Durability*100)
}

func (r Rod) level() string {
	v := [...]string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X"}
	return v[int(r.Quality*float64(len(v)))]
}

func NewRod() *Rod {
	f := &Rod{
		Quality:    0.5 + 0.2*rand.NormFloat64(),
		Durability: 0.8 + 0.2*rand.Float64(),
	}
	if f.Quality < 0 || f.Quality > 1 {
		return NewRod()
	}
	return f
}
