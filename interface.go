package autotile

import (
	"github.com/voidshard/tile"
)

// Outline tells us what we can find at given map locations
type Outline interface {
	LandAt(x, y int) LandData
}

// LandData represents information about the natural world at a given location
type LandData interface {
	// Asks if the given tile is one of these things.
	// One of these should return true
	IsLand() bool
	IsWater() bool
	IsMolten() bool
	IsNull() bool

	IsRoad() bool

	// statistics that affect the kinds of things we draw.
	// The units aren't important, other than they make sense when compared to
	// the autotiler Config biome settings.
	Height() int
	Rainfall() int
	Temperature() int

	// actual tiles we can place
	Tiles() *LandTiles

	// user defined tags for this (x,y)
	Tags() []string
}

// ObjectBin is something that can decide what object should be placed at
// a given location on a map.
type ObjectBin interface {
	// Choose returns the object ID placed, the object itself and/or an error
	Choose(t tile.Tileable, x, y, z int) (string, *tile.Map, error)
}
