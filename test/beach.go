package main

import ()

// defines a lovely beach scene
var beach *testOutline

func init() {
	midTemp := (cfg.VegetationMaxTemp-cfg.VegetationMinTemp)/2 + cfg.VegetationMinTemp
	cliffl := cfg.CliffLevel

	beach = newOutline([][]*Area{
		[]*Area{ // 0
			&Area{},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1, River: true},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
		},

		[]*Area{ // 1
			&Area{},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1, River: true},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
		},
		[]*Area{ // 2
			&Area{},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1, River: true},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
		},
		[]*Area{ // 3
			&Area{},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1, River: true},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl},
			&Area{height: cliffl},
			&Area{height: cliffl},
		},
		[]*Area{ // 4
			&Area{},
			&Area{},
			&Area{River: true},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{height: cliffl - 1},
			&Area{height: cliffl - 1},
			&Area{height: cliffl - 1},
		},
		[]*Area{ // 5
			&Area{},
			&Area{},
			&Area{River: true},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
		},
		[]*Area{ // 6
			&Area{},
			&Area{River: true},
			&Area{River: true},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
		},
		[]*Area{ // 7
			&Area{River: true},
			&Area{River: true},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{Sea: true},
			&Area{Sea: true},
		},
		[]*Area{ // 8
			&Area{},
			&Area{River: true},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{Sea: true},
			&Area{Sea: true},
			&Area{Sea: true},
			&Area{Sea: true},
		},
		[]*Area{ // 9
			&Area{Sea: true},
			&Area{Sea: true, River: true},
			&Area{Sea: true},
			&Area{Sea: true},
			&Area{Sea: true},
			&Area{Sea: true},
			&Area{Sea: true},
			&Area{Sea: true},
			&Area{Sea: true},
			&Area{Sea: true},
		},
	})
	beach.DefaultTemp = midTemp
	beach.SetScale(3)
}
