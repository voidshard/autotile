package main

var cliffs *testOutline

func init() {
	midTemp := (cfg.VegetationMaxTemp-cfg.VegetationMinTemp)/2 + cfg.VegetationMinTemp
	cliffl := cfg.CliffLevel

	cliffs = newOutline([][]*Area{

		[]*Area{ // 0
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{Road: true},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
		},
		[]*Area{ // 0
			&Area{},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 1},
			&Area{height: cliffl + 10},
			&Area{height: cliffl + 10},
			&Area{Road: true, height: cliffl + 10},
			&Area{height: cliffl + 10},
			&Area{height: cliffl + 10},
			&Area{height: cliffl + 1},
			&Area{},
		},
		[]*Area{ // 0
			&Area{Road: true},
			&Area{height: cliffl + 1, Road: true},
			&Area{height: cliffl + 1, Road: true},
			&Area{height: cliffl + 10, Road: true},
			&Area{height: cliffl + 10, Road: true},
			&Area{height: cliffl + 10, Road: true},
			&Area{height: cliffl + 10, Road: true},
			&Area{height: cliffl + 1, Road: true},
			&Area{height: cliffl + 1, Road: true},
			&Area{Road: true},
		},
		[]*Area{ // 0
			&Area{River: true},
			&Area{height: cliffl + 1, River: true},
			&Area{height: cliffl + 1, River: true},
			&Area{height: cliffl + 10, River: true},
			&Area{height: cliffl + 10, River: true},
			&Area{height: cliffl + 10, Road: true, River: true},
			&Area{height: cliffl + 10, River: true},
			&Area{height: cliffl + 10, River: true},
			&Area{height: cliffl + 1, River: true},
			&Area{River: true},
		},
		[]*Area{ // 0
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
			&Area{Road: true},
			&Area{},
			&Area{},
			&Area{},
			&Area{},
		},
	})
	cliffs.DefaultTemp = midTemp
	cliffs.SetScale(3)
}
