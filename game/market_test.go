package game

import "testing"

func TestRandomProduct(t *testing.T) {
	const sample = 100000
	prices := map[int]int{}
	for i := 0; i < sample; i++ {
		p := randomProduct()
		if p.Price < 0 {
			t.Errorf("negative product price %v", p.Price)
		}
		prices[p.Price]++
	}
	for price, total := range prices {
		const threshold = 1000
		if total > threshold {
			t.Errorf("price %v encountered %v times", price, total)
		}
	}
}
