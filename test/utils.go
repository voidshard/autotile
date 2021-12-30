package main

import (
	"image"

	"github.com/voidshard/autotile"
)

var cfg = &autotile.Config{
	StepSize: 10,
	MapSize:  30,
	TileSize: 32,
	Params: &autotile.WorldParams{
		BeachWidth:        3,
		VegetationMaxTemp: 45,
		VegetationMinTemp: -5,
		MountainLevel:     240,
		CliffLevel:        170,
	},
	Routines: 10,
}

var land = &autotile.Land{
	Grass: &autotile.BasicLand{
		Full:       []string{"grass.full.01.png"},
		Transition: []string{"grass.trans.01.png"},
	},
	Sand: &autotile.BasicLand{
		Full:       []string{"sand.full.01.png"},
		Transition: []string{"sand.trans.01.png"},
	},
	Dirt: &autotile.BasicLand{
		Full:       []string{"dirt.full.01.png"},
		Transition: []string{"dirt.trans.01.png"},
	},
	Snow: &autotile.BasicLand{
		Full:       []string{"snow.full.01.png"},
		Transition: []string{"snow.trans.01.png"},
	},
	Rock: &autotile.BasicLand{
		Full:       []string{"rock.full.01.png"},
		Transition: []string{"rock.trans.01.png"},
	},
	Water: &autotile.ComplexLand{
		Full:                  []string{"water.full.01.png"},
		NorthHalf:             []string{"water.n.01.png"},
		EastHalf:              []string{"water.e.01.png"},
		SouthHalf:             []string{"water.s.01.png"},
		WestHalf:              []string{"water.w.01.png"},
		QuarterNorthEast:      []string{"water.1q_ne.01.png"},
		QuarterSouthEast:      []string{"water.1q_se.01.png"},
		QuarterSouthWest:      []string{"water.1q_sw.01.png"},
		QuarterNorthWest:      []string{"water.1q_nw.01.png"},
		ThreeQuarterNorthEast: []string{"water.3q_ne.01.png"},
		ThreeQuarterSouthEast: []string{"water.3q_se.01.png"},
		ThreeQuarterSouthWest: []string{"water.3q_sw.01.png"},
		ThreeQuarterNorthWest: []string{"water.3q_nw.01.png"},
	},
	Cliff: &autotile.Cliff{
		NorthHalf:             []string{"cliff.n.01.png"},
		EastHalf:              []string{"cliff.e.01.png"},
		SouthHalf:             []string{"cliff.s.01.png"},
		WestHalf:              []string{"cliff.w.01.png"},
		QuarterNorthEast:      []string{"cliff.1q_ne.01.png"},
		QuarterSouthEast:      []string{"cliff.1q_se.01.png"},
		QuarterSouthWest:      []string{"cliff.1q_sw.01.png"},
		QuarterNorthWest:      []string{"cliff.1q_nw.01.png"},
		NorthHalfBase:         []string{"cliff.n_base.01.png"},
		QuarterNorthEastBase:  []string{"cliff.1q_ne_base.01.png"},
		QuarterNorthWestBase:  []string{"cliff.1q_nw_base.01.png"},
		ThreeQuarterNorthEast: []string{"cliff.3q_ne.01.png"},
		ThreeQuarterNorthWest: []string{"cliff.3q_nw.01.png"},
	},
	Waterfall: &autotile.Waterfall{
		NS: &autotile.WaterfallNorthSouth{
			LeftTop:     []string{"waterfall.ns.01.0.0.0.png"},
			LeftCentre:  []string{"waterfall.ns.01.0.1.0.png"},
			LeftBottom:  []string{"waterfall.ns.01.0.2.0.png"},
			MidTop:      []string{"waterfall.ns.01.1.0.0.png"},
			MidCentre:   []string{"waterfall.ns.01.1.1.0.png"},
			MidBottom:   []string{"waterfall.ns.01.1.2.0.png"},
			RightTop:    []string{"waterfall.ns.01.2.0.0.png"},
			RightCentre: []string{"waterfall.ns.01.2.1.0.png"},
			RightBottom: []string{"waterfall.ns.01.2.2.0.png"},
		},
		SN: &autotile.WaterfallSouthNorth{
			MidTop: []string{"waterfall.sn.top.01.0.0.0.png"},
		},
	},
}

type testOutline struct {
	r           image.Rectangle
	grid        [][]*autotile.Area
	DefaultTemp int
}

func newOutline(grid [][]*autotile.Area) *testOutline {
	max := len(grid)
	for _, row := range grid {
		if len(row) > max {
			max = len(row)
		}
	}
	return &testOutline{
		r:    image.Rect(0, 0, max, max),
		grid: grid,
	}
}

func (o *testOutline) Bounds() image.Rectangle {
	return o.r
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (o *testOutline) At(x, y int) *autotile.Area {
	if y >= 0 && y < len(o.grid) {
		if x >= 0 && x < len(o.grid[y]) {
			a := o.grid[y][x]
			a.X = x
			a.Y = y
			if a.Temperature == 0 {
				a.Temperature = o.DefaultTemp
			}
			a.Land = land
			return a
		}
	}

	validy := min(max(y, 0), len(o.grid)-1)
	validx := min(max(x, 0), len(o.grid[validy])-1)

	return o.At(validx, validy)
}
