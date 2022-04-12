package main

import (
	"fmt"

	"github.com/voidshard/autotile"
	"github.com/voidshard/autotile/pkg/estate"
)

var (
	house = estate.NewObject("house.01.tmx")
	corn  = estate.NewObject("wildvege.05.tmx").SetMax(12).SetMin(0)
	tree  = estate.NewObject("tree.01.tmx").SetMax(3).SetMin(0)
	fence = &autotile.Tileset{
		NorthHalf:        []string{"fence.n.01.0.0.0.png"},
		EastHalf:         []string{"fence.e.01.0.0.0.png"},
		WestHalf:         []string{"fence.w.01.0.0.0.png"},
		SouthHalf:        []string{"fence.s.01.0.0.0.png"},
		QuarterNorthEast: []string{"fence.1q_ne.01.0.0.0.png"},
		QuarterNorthWest: []string{"fence.1q_nw.01.0.0.0.png"},
		QuarterSouthEast: []string{"fence.1q_se.01.0.0.0.png"},
		QuarterSouthWest: []string{"fence.1q_sw.01.0.0.0.png"},
	}
	dirt = &autotile.Tileset{
		Full:                  []string{"dirt.full.01.0.0.0.png"},
		NorthHalf:             []string{"dirt.n.01.0.0.0.png"},
		EastHalf:              []string{"dirt.e.01.0.0.0.png"},
		SouthHalf:             []string{"dirt.s.01.0.0.0.png"},
		WestHalf:              []string{"dirt.w.01.0.0.0.png"},
		QuarterNorthEast:      []string{"dirt.1q_ne.01.0.0.0.png"},
		QuarterSouthEast:      []string{"dirt.1q_se.01.0.0.0.png"},
		QuarterNorthWest:      []string{"dirt.1q_nw.01.0.0.0.png"},
		QuarterSouthWest:      []string{"dirt.1q_sw.01.0.0.0.png"},
		ThreeQuarterNorthEast: []string{"dirt.3q_ne.01.0.0.0.png"},
		ThreeQuarterNorthWest: []string{"dirt.3q_nw.01.0.0.0.png"},
		ThreeQuarterSouthEast: []string{"dirt.3q_se.01.0.0.0.png"},
		ThreeQuarterSouthWest: []string{"dirt.3q_sw.01.0.0.0.png"},
	}
)

func main() {
	fmt.Println("test estate")

	cfg := &estate.Config{
		// set a rng seed
		Seed: 123456789,
		// tell us how to load tile.Map objects from strings
		Loader: autotile.NewFileLoader("test/tobs"),
		Set: &estate.Set{
			// insert empty space ~50% of space used for objects in this set
			EmptyPercentage: 0.5,
			// fill with house & tree(s)
			Objects: []*estate.Object{house, tree},
			// leave empty tile on left, right, bottom (of all non-empty space)
			PadLeft:   1,
			PadRight:  1,
			PadBottom: 1,
			// surround this set with fence
			Fence: fence,
			// leave hole in fence (since we don't set Gate)
			GateLocation: estate.CentreBottom,
			Sets: []*estate.Set{
				// garden is a sub set
				&estate.Set{
					// fill with corn
					Objects: []*estate.Object{corn},
					// and fence
					Fence: fence,
					// again leave a gap at the bottom for a gate
					GateLocation: estate.CentreBottom,
					// and fill the garden with dirt
					Base: dirt,
				},
			},
		},
	}

	e, err := estate.Build(cfg)
	if err != nil {
		panic(err)
	}

	m, err := e.Map(16, 16)
	if err != nil {
		panic(err)
	}

	err = m.WriteFile("estate.tmx")
	if err != nil {
		panic(err)
	}
}
