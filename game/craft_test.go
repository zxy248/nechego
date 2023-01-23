package game

import (
	"math"
	"nechego/fishing"
	"nechego/item"
	"testing"
)

func TestCraftRod(t *testing.T) {
	table := []struct {
		recipe []*fishing.Rod
		want   *fishing.Rod
	}{
		{
			[]*fishing.Rod{
				{Quality: 0.5, Durability: 0.5},
				{Quality: 0.5, Durability: 0.5},
			},
			&fishing.Rod{Quality: 0.6, Durability: 1.0},
		},
		{
			[]*fishing.Rod{
				{Quality: 0.0, Durability: 0.3},
				{Quality: 1.0, Durability: 0.4},
			},
			&fishing.Rod{Quality: 0.6, Durability: 0.8},
		},
		{
			[]*fishing.Rod{
				{Quality: 1.0, Durability: 0.0},
				{Quality: 1.0, Durability: 0.0},
			},
			&fishing.Rod{Quality: 1.1, Durability: 0.1},
		},
		{
			[]*fishing.Rod{
				{Quality: 0.0, Durability: 1.0},
				{Quality: 0.0, Durability: 1.0},
			},
			&fishing.Rod{Quality: 0.1, Durability: 1.0},
		},
	}
	for _, x := range table {
		result, _ := craftRod([]*item.Item{{Value: x.recipe[0]}, {Value: x.recipe[1]}})
		rod := result[0].Value.(*fishing.Rod)
		if math.Abs(rod.Durability-x.want.Durability) > 0.001 {
			t.Errorf("got %v, want %v", rod.Durability, x.want.Durability)
		}
		if rod.Quality != x.want.Quality {
			t.Errorf("got %v, want %v", rod.Quality, x.want.Quality)
		}
	}
}
