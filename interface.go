package autotile

import (
	"github.com/voidshard/tile"
)

//
type Outline interface {
	LandAt(x, y int) LandData
}

//
type LandData interface {
	// Asks if the given tile is one of these things.
	// One of these should return true
	IsLand() bool
	IsWater() bool
	IsMolten() bool
	IsNull() bool

	//
	IsRoad() bool

	// statistics that affect the kinds of things we draw
	Height() int
	Rainfall() int
	Temperature() int

	// actual tiles we can place
	Tiles() *LandTiles

	// user defined tags for this (x,y)
	Tags() []string
}

//
type ObjectBin interface {
	Choose(t tile.Tileable, x, y, z int) (string, *tile.Map, error)
}
