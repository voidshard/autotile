package main

import (
	"image"

	"github.com/voidshard/autotile"
)

var cfg = &autotile.Config{
	BeachWidth:        2,
	VegetationMaxTemp: 45,
	VegetationMinTemp: -5,
	MountainLevel:     240,
	CliffLevel:        170,
}

var land = &autotile.LandTiles{
	Grass: &autotile.Tileset{
		Full: []string{"grass.full.01.0.0.0.png"},
	},
	Dirt: &autotile.Tileset{
		Full: []string{"dirt.full.01.0.0.0.png"},
	},
	Sand: &autotile.Tileset{
		Full: []string{"sand.full.01.0.0.0.png"},
	},
	Water: &autotile.Tileset{
		Full:                  []string{"river.full.01.0.0.0.png"},
		NorthHalf:             []string{"river.n.01.0.0.0.png"},
		EastHalf:              []string{"river.e.01.0.0.0.png"},
		SouthHalf:             []string{"river.s.01.0.0.0.png"},
		WestHalf:              []string{"river.w.01.0.0.0.png"},
		QuarterNorthEast:      []string{"river.1q_ne.01.0.0.0.png"},
		QuarterSouthEast:      []string{"river.1q_se.01.0.0.0.png"},
		QuarterSouthWest:      []string{"river.1q_sw.01.0.0.0.png"},
		QuarterNorthWest:      []string{"river.1q_nw.01.0.0.0.png"},
		ThreeQuarterNorthEast: []string{"river.3q_ne.01.0.0.0.png"},
		ThreeQuarterSouthEast: []string{"river.3q_se.01.0.0.0.png"},
		ThreeQuarterSouthWest: []string{"river.3q_sw.01.0.0.0.png"},
		ThreeQuarterNorthWest: []string{"river.3q_nw.01.0.0.0.png"},
	},
	Bridge: &autotile.Tileset{
		Full:                  []string{"bridge.plank.full.01.0.0.0.png"},
		NorthHalf:             []string{"bridge.plank.n.01.0.0.0.png"},
		EastHalf:              []string{"bridge.plank.e.01.0.0.0.png"},
		SouthHalf:             []string{"bridge.plank.s.01.0.0.0.png"},
		WestHalf:              []string{"bridge.plank.w.01.0.0.0.png"},
		QuarterNorthEast:      []string{"bridge.plank.1q_ne.01.0.0.0.png"},
		QuarterSouthEast:      []string{"bridge.plank.1q_se.01.0.0.0.png"},
		QuarterSouthWest:      []string{"bridge.plank.1q_sw.01.0.0.0.png"},
		QuarterNorthWest:      []string{"bridge.plank.1q_nw.01.0.0.0.png"},
		ThreeQuarterNorthEast: []string{"bridge.plank.3q_ne.01.0.0.0.png"},
		ThreeQuarterNorthWest: []string{"bridge.plank.3q_nw.01.0.0.0.png"},
		ThreeQuarterSouthEast: []string{"bridge.plank.3q_se.01.0.0.0.png"},
		ThreeQuarterSouthWest: []string{"bridge.plank.3q_sw.01.0.0.0.png"},
	},
	Cliff: &autotile.Tileset{
		NorthHalf:             []string{"cliffs.n.01.0.0.0.png"},
		EastHalf:              []string{"cliffs.e.01.0.0.0.png"},
		SouthHalf:             []string{"cliffs.s.01.0.0.0.png"},
		WestHalf:              []string{"cliffs.w.01.0.0.0.png"},
		QuarterNorthEast:      []string{"cliffs.1q_ne.01.0.0.0.png"},
		QuarterSouthEast:      []string{"cliffs.1q_se.01.0.0.0.png"},
		QuarterSouthWest:      []string{"cliffs.1q_sw.01.0.0.0.png"},
		QuarterNorthWest:      []string{"cliffs.1q_nw.01.0.0.0.png"},
		ThreeQuarterNorthEast: []string{"cliffs.3q_ne.01.0.0.0.png"},
		ThreeQuarterNorthWest: []string{"cliffs.3q_nw.01.0.0.0.png"},
		ThreeQuarterSouthEast: []string{"cliffs.3q_se.01.0.0.0.png"},
		ThreeQuarterSouthWest: []string{"cliffs.3q_sw.01.0.0.0.png"},
	},
	Road: &autotile.Tileset{
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
	},
	StairsNorthSouth: &autotile.Tileset{
		Full:             []string{"stairs.ns.full.01.0.0.0.png"},
		NorthHalf:        []string{"stairs.ns.n.01.0.0.0.png"},
		EastHalf:         []string{"stairs.ns.e.01.0.0.0.png"},
		SouthHalf:        []string{"stairs.ns.s.01.0.0.0.png"},
		WestHalf:         []string{"stairs.ns.w.01.0.0.0.png"},
		QuarterNorthEast: []string{"stairs.ns.1q_ne.01.0.0.0.png"},
		QuarterSouthEast: []string{"stairs.ns.1q_se.01.0.0.0.png"},
		QuarterNorthWest: []string{"stairs.ns.1q_nw.01.0.0.0.png"},
		QuarterSouthWest: []string{"stairs.ns.1q_sw.01.0.0.0.png"},
	},
	StairsWestEast: &autotile.Tileset{
		Full:             []string{"stairs.we.full.01.0.0.0.png"},
		NorthHalf:        []string{"stairs.we.n.01.0.0.0.png"},
		EastHalf:         []string{"stairs.we.e.01.0.0.0.png"},
		SouthHalf:        []string{"stairs.we.s.01.0.0.0.png"},
		WestHalf:         []string{"stairs.we.w.01.0.0.0.png"},
		QuarterNorthEast: []string{"stairs.we.1q_ne.01.0.0.0.png"},
		QuarterSouthEast: []string{"stairs.we.1q_se.01.0.0.0.png"},
		QuarterNorthWest: []string{"stairs.we.1q_nw.01.0.0.0.png"},
		QuarterSouthWest: []string{"stairs.we.1q_sw.01.0.0.0.png"},
	},
	StairsSouthNorth: &autotile.Tileset{
		Full:             []string{"stairs.ns.full.01.0.0.0.png"},
		NorthHalf:        []string{"stairs.ns.n.01.0.0.0.png"},
		EastHalf:         []string{"stairs.ns.e.01.0.0.0.png"},
		SouthHalf:        []string{"stairs.ns.s.01.0.0.0.png"},
		WestHalf:         []string{"stairs.ns.w.01.0.0.0.png"},
		QuarterNorthEast: []string{"stairs.ns.1q_ne.01.0.0.0.png"},
		QuarterSouthEast: []string{"stairs.ns.1q_se.01.0.0.0.png"},
		QuarterNorthWest: []string{"stairs.ns.1q_nw.01.0.0.0.png"},
		QuarterSouthWest: []string{"stairs.ns.1q_sw.01.0.0.0.png"},
	},
	StairsEastWest: &autotile.Tileset{
		Full:             []string{"stairs.ew.full.01.0.0.0.png"},
		NorthHalf:        []string{"stairs.ew.n.01.0.0.0.png"},
		EastHalf:         []string{"stairs.ew.e.01.0.0.0.png"},
		SouthHalf:        []string{"stairs.ew.s.01.0.0.0.png"},
		WestHalf:         []string{"stairs.ew.w.01.0.0.0.png"},
		QuarterNorthEast: []string{"stairs.ew.1q_ne.01.0.0.0.png"},
		QuarterSouthEast: []string{"stairs.ew.1q_se.01.0.0.0.png"},
		QuarterNorthWest: []string{"stairs.ew.1q_nw.01.0.0.0.png"},
		QuarterSouthWest: []string{"stairs.ew.1q_sw.01.0.0.0.png"},
	},
	WaterfallNorthSouth: &autotile.Tileset{
		Full:             []string{"waterfall.ns.full.01.0.0.0.png"},
		NorthHalf:        []string{"waterfall.ns.n.01.0.0.0.png"},
		EastHalf:         []string{"waterfall.ns.e.01.0.0.0.png"},
		SouthHalf:        []string{"waterfall.ns.s.01.0.0.0.png"},
		WestHalf:         []string{"waterfall.ns.w.01.0.0.0.png"},
		QuarterNorthEast: []string{"waterfall.ns.1q_ne.01.0.0.0.png"},
		QuarterSouthEast: []string{"waterfall.ns.1q_se.01.0.0.0.png"},
		QuarterNorthWest: []string{"waterfall.ns.1q_nw.01.0.0.0.png"},
		QuarterSouthWest: []string{"waterfall.ns.1q_sw.01.0.0.0.png"},
	},
	WaterfallSouthNorth: &autotile.Tileset{
		SouthHalf:        []string{"waterfall.sn.s.01.0.0.0.png"},
		QuarterSouthEast: []string{"waterfall.sn.1q_se.01.0.0.0.png"},
		QuarterSouthWest: []string{"waterfall.sn.1q_sw.01.0.0.0.png"},
	},
	WaterfallEastWest: &autotile.Tileset{
		QuarterNorthEast: []string{"waterfall.ew.1q_ne.01.0.0.0.png"},
		QuarterNorthWest: []string{"waterfall.ew.1q_nw.01.0.0.0.png"},
		QuarterSouthEast: []string{"waterfall.ew.1q_se.01.0.0.0.png"},
		QuarterSouthWest: []string{"waterfall.ew.1q_sw.01.0.0.0.png"},
		EastHalf:         []string{"waterfall.ew.e.01.0.0.0.png"},
		Full:             []string{"waterfall.ew.full.01.0.0.0.png"},
		NorthHalf:        []string{"waterfall.ew.n.01.0.0.0.png"},
		SouthHalf:        []string{"waterfall.ew.s.01.0.0.0.png"},
		WestHalf:         []string{"waterfall.ew.w.01.0.0.0.png"},
	},
	WaterfallWestEast: &autotile.Tileset{
		QuarterNorthEast: []string{"waterfall.we.1q_ne.01.0.0.0.png"},
		QuarterNorthWest: []string{"waterfall.we.1q_nw.01.0.0.0.png"},
		QuarterSouthEast: []string{"waterfall.we.1q_se.01.0.0.0.png"},
		QuarterSouthWest: []string{"waterfall.we.1q_sw.01.0.0.0.png"},
		EastHalf:         []string{"waterfall.we.e.01.0.0.0.png"},
		Full:             []string{"waterfall.we.full.01.0.0.0.png"},
		NorthHalf:        []string{"waterfall.we.n.01.0.0.0.png"},
		SouthHalf:        []string{"waterfall.we.s.01.0.0.0.png"},
		WestHalf:         []string{"waterfall.we.w.01.0.0.0.png"},
	},
}

type testOutline struct {
	grid [][]*Area
	max  int

	scale       int
	DefaultTemp int
}

func newOutline(grid [][]*Area) *testOutline {
	max := len(grid)
	for _, row := range grid {
		if len(row) > max {
			max = len(row)
		}
	}
	return &testOutline{
		grid:  grid,
		max:   max,
		scale: 1,
	}
}

func (o *testOutline) SetScale(i int) {
	if i <= 0 {
		o.SetScale(1)
	}
	o.scale = i
}

func (o *testOutline) Bounds() image.Rectangle {
	return image.Rect(0, 0, o.max*o.scale, o.max*o.scale)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (o *testOutline) LandAt(inx, iny int) autotile.LandData {
	inx = maxInt(minInt(inx, o.max*o.scale-1), 0)
	iny = maxInt(minInt(iny, o.max*o.scale-1), 0)

	x := inx / o.scale
	y := iny / o.scale

	if y >= 0 && y < len(o.grid) {
		if x >= 0 && x < len(o.grid[y]) {
			result := o.grid[y][x]
			result.o = o
			return result
		}
	}

	validy := minInt(maxInt(y, 0), len(o.grid)-1)
	validx := minInt(maxInt(x, 0), len(o.grid[validy])-1)

	return o.LandAt(validx, validy)
}

type Area struct {
	Sea    bool
	River  bool
	Road   bool
	height int
	o      *testOutline
}

func (a *Area) IsLand() bool               { return !a.IsWater() }
func (a *Area) IsWater() bool              { return a.Sea || a.River }
func (a *Area) IsMolten() bool             { return false }
func (a *Area) IsRoad() bool               { return a.Road }
func (a *Area) IsNull() bool               { return false }
func (a *Area) Height() int                { return a.height }
func (a *Area) Rainfall() int              { return 150 }
func (a *Area) Temperature() int           { return a.o.DefaultTemp }
func (a *Area) Tiles() *autotile.LandTiles { return land }
func (a *Area) Tags() []string             { return nil }
