package main

import (
	"github.com/voidshard/autotile"
)

// defines a lovely beach scene
var beach *testOutline

func init() {
	midTemp := (cfg.Params.VegetationMaxTemp-cfg.Params.VegetationMinTemp)/2 + cfg.Params.VegetationMinTemp
	cliffl := cfg.Params.CliffLevel

	beach = newOutline([][]*autotile.Area{
		[]*autotile.Area{ // 0
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl + 1, River: true},
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl + 1},
		},
		[]*autotile.Area{ // 1
			&autotile.Area{},
			&autotile.Area{Height: cliffl + 1, River: true},
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl + 1},
			&autotile.Area{Height: cliffl},
			&autotile.Area{Height: cliffl},
			&autotile.Area{Height: cliffl},
		},
		[]*autotile.Area{ // 2
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{River: true},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{Height: cliffl - 1},
			&autotile.Area{Height: cliffl - 1},
			&autotile.Area{Height: cliffl - 1},
		},
		[]*autotile.Area{ // 3
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{River: true},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
		},
		[]*autotile.Area{ // 4
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{River: true},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
		},
		[]*autotile.Area{ // 5
			&autotile.Area{},
			&autotile.Area{River: true},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
		},
		[]*autotile.Area{ // 6
			&autotile.Area{},
			&autotile.Area{River: true},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
		},
		[]*autotile.Area{ // 7
			&autotile.Area{River: true},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{Sea: true},
			&autotile.Area{Sea: true},
		},
		[]*autotile.Area{ // 8
			&autotile.Area{},
			&autotile.Area{River: true},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{},
			&autotile.Area{Sea: true},
			&autotile.Area{Sea: true},
			&autotile.Area{Sea: true},
			&autotile.Area{Sea: true},
		},
		[]*autotile.Area{ // 9
			&autotile.Area{Sea: true},
			&autotile.Area{Sea: true, River: true},
			&autotile.Area{Sea: true},
			&autotile.Area{Sea: true},
			&autotile.Area{Sea: true},
			&autotile.Area{Sea: true},
			&autotile.Area{Sea: true},
			&autotile.Area{Sea: true},
			&autotile.Area{Sea: true},
			&autotile.Area{Sea: true},
		},
	})
	beach.DefaultTemp = midTemp
}
