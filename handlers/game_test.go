package handlers

import (
	"math"
	"nechego/game"
	"testing"
)

// Telegram IDs are usually large numbers.
const tuidOffset = 250_000_000

func TestFishSuccessChance(t *testing.T) {
	const sample = 100000
	total := float64(0)
	for i := 0; i < sample; i++ {
		u := &game.User{TUID: tuidOffset + int64(i)}
		chance := fishSuccessChance(u)
		total += chance
		const min, max = -0.02, 1.02
		if chance < min || chance > max {
			t.Errorf("chance = %v, want in [%v, %v]", chance, min, max)
		}
	}
	t.Run("average", func(t *testing.T) {
		const want = 0.5
		avg := total / sample
		delta := math.Abs(avg - want)
		if delta > 0.01 {
			t.Errorf("delta = %v, want epsilon", delta)
		}
	})
}
