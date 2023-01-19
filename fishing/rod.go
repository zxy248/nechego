package fishing

import (
	"fmt"
	"math/rand"
)

var (
	levels = [...]string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X"}
	greeks = [...]string{"α", "β", "γ", "δ", "ε", "ζ", "η", "θ", "ι", "κ", "λ",
		"μ", "ν", "ξ", "ο", "π", "ρ", "σ", "τ", "υ", "φ", "χ", "ψ", "ω"}
)

type Rod struct {
	Quality    float64 // [0, 1] (initially)
	Durability float64 // [0, 1]
}

func (r Rod) String() string {
	return fmt.Sprintf("🎣 Удочка (%s, %.f%%)", r.level(), r.Durability*100)
}

func (r Rod) level() string {
	switch q := r.Quality; {
	case q >= 0 && q < 1:
		return levels[int(q*float64(len(levels)))]
	case q >= 1 && q < 3.4:
		q = (q - 1) / (0.1 * float64(len(greeks)))
		return greeks[int(q*float64(len(greeks)))]
	default:
		panic(fmt.Errorf("unexpected quality %v", r.Quality))
	}
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
