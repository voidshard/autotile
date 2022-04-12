package estate

import (
	"github.com/voidshard/tile"

	"github.com/voidshard/autotile"
)

type Location string

const (
	// CentreBottom indicates position on a rectangle
	CentreBottom = "centre-bottom"

	// None indicates something ought not be placed
	None = ""
)

// Set is some objects that are placed together.
// Each item from `Objects` will be placed in a rectangle.
// If given, the entire region will have base tiles set from `Base`
// If given, the region will be surrounded by `Fence`
type Set struct {
	// Objects (tobs) that will be placed
	Objects []*Object

	// Pad left side of objects with empty tiles
	PadLeft int

	// PadRight side of objects with empty tiles
	PadRight int

	// PadTop side of objects with empty tiles
	PadTop int

	// PadBottom side of objects with empty tiles
	PadBottom int

	// Base tiles to go under the entire region
	Base      *autotile.Tileset
	BaseProps *tile.Properties

	// Fence to surround the region with
	Fence      *autotile.Tileset
	FenceProps *tile.Properties

	// Gate image (or "" if you simply wish a blank space)
	// nb. implies fence
	Gate string

	// Indicates where in a fence a Gate should be set.
	// nb. implies fence
	GateLocation Location

	// Sets that this set contains
	Sets []*Set

	// rough % of tiles we'll add just to have empty space
	EmptyPercentage float64
}
