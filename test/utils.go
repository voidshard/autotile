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
	height int
	o      *testOutline
}

func (a *Area) IsLand() bool               { return !a.IsWater() }
func (a *Area) IsWater() bool              { return a.Sea || a.River }
func (a *Area) IsMolten() bool             { return false }
func (a *Area) IsRoad() bool               { return false }
func (a *Area) IsNull() bool               { return false }
func (a *Area) Height() int                { return a.height }
func (a *Area) Rainfall() int              { return 150 }
func (a *Area) Temperature() int           { return a.o.DefaultTemp }
func (a *Area) Tiles() *autotile.LandTiles { return land }
func (a *Area) Tags() []string             { return nil }
