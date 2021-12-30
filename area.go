package autotile

import (
	"sort"
)

var (
	// helpful corners
	cornerNE = []Heading{North, NorthEast, East}
	cornerNW = []Heading{West, NorthWest, North}
	cornerSE = []Heading{East, SouthEast, South}
	cornerSW = []Heading{South, SouthWest, West}

	edgeNE = []Heading{
		NorthWest,
		North,
		NorthEast,
		East,
		SouthEast,
	}
	edgeNW = []Heading{
		SouthWest,
		West,
		NorthWest,
		North,
		NorthEast,
	}
	edgeSE = []Heading{
		NorthEast,
		East,
		SouthEast,
		South,
		SouthWest,
	}
	edgeSW = []Heading{
		SouthEast,
		South,
		SouthWest,
		West,
		NorthWest,
	}
)

// Area represents a single pixel section on a map from which
// we can generate tiles.
type Area struct {
	// the X co-ord of this area
	X int

	// the Y co-ord of this area
	Y int

	// land tiles applicable to this area.
	Land *Land

	// height in units - actual units isn't important.
	// Differences in height over our config.CliffLevel is used to decide
	// where cliffs go. We also use this to transition to rock landscapes over
	// config.MountainLevel
	Height int

	// Temperature in degrees. We use this along with our config values
	// VegetationMinTemp, VegetationMaxTemp, SnowLevel to determine appropriate
	// terrain types to use (ie. snow, grass, dirt, rock) etc.
	Temperature int

	// Sea true if this pixel is in the sea
	Sea bool

	// River true if this pixel is in a river
	River bool

	// Swamp true if this pixel is swamp water
	Swamp bool

	// Lava true if there is molten rock here
	Lava bool

	// Road true if we should be depicting a road/path/street
	Road bool

	// Tags are custom tags added to this tile(s) for use with an
	// objectbin later on
	Tags []string

	// internally used to track this tiles' relation to another
	heading Heading
}

// isWater returns if we consider this area underwater
func (a *Area) isWater() bool {
	return a.Sea || a.River || a.Swamp
}

// setHeading sets heading to `in` and returns the area
func (a *Area) setHeading(in Heading) *Area {
	a.heading = in
	return a
}

// headings returns the `heading` values of the given areas
func headings(in []*Area) []Heading {
	ls := []Heading{}
	for _, h := range in {
		ls = append(ls, h.heading)
	}
	sort.Slice(ls, func(i, j int) bool { return int(ls[i]) < int(ls[j]) })
	return ls
}

// includes returns if list a contains list b (in order)
func includes(a []Heading, b ...Heading) bool {
	if a == nil || b == nil || len(b) > len(a) || len(b) == 0 {
		return false
	}

	if len(b) > 1 {
		a = append(a, a[0:len(b)-1]...) // allows us to detect when a contains b but wrapped around
	}
	for i := 0; i <= len(a)-len(b); i++ {
		reject := false
		for j := range b {
			if a[i+j] != b[j] {
				reject = true
				break
			}
		}
		if !reject {
			return true
		}
	}
	return false
}

// nearby is a set of areas near the given tile
type nearby struct {
	North     *Area
	NorthEast *Area
	East      *Area
	SouthEast *Area
	South     *Area
	SouthWest *Area
	West      *Area
	NorthWest *Area
}

// all returns all tiles nearby in 8 cardinal directions
func (n *nearby) all() []*Area {
	return []*Area{
		n.North,
		n.NorthEast,
		n.East,
		n.SouthEast,
		n.South,
		n.SouthWest,
		n.West,
		n.NorthWest,
	}
}

// Lower returns tiles nearby lower than the given value (in height)
func (n *nearby) Lower(h int) []*Area {
	land := []*Area{}
	for _, t := range n.all() {
		if t.Height < h {
			land = append(land, t)
		}
	}
	return land
}

// Land returns all tiles nearby that are not water
func (n *nearby) Land() []*Area {
	land := []*Area{}
	for _, t := range n.all() {
		if t.isWater() {
			continue
		}
		land = append(land, t)
	}
	return land
}

// Water returns all tiles nearby that are water
func (n *nearby) Water() []*Area {
	water := []*Area{}
	for _, t := range n.all() {
		if t.isWater() {
			water = append(water, t)
		}
	}
	return water
}

// cardinals returns information on surrounding tiles
func (a *Autotiler) cardinals(o Outline, mx, my, tx, ty int) *nearby {
	return &nearby{
		a.AtMapCoord(o, mx, my, tx, ty-1).setHeading(North),
		a.AtMapCoord(o, mx, my, tx+1, ty-1).setHeading(NorthEast),
		a.AtMapCoord(o, mx, my, tx+1, ty).setHeading(East),
		a.AtMapCoord(o, mx, my, tx+1, ty+1).setHeading(SouthEast),
		a.AtMapCoord(o, mx, my, tx, ty+1).setHeading(South),
		a.AtMapCoord(o, mx, my, tx-1, ty+1).setHeading(SouthWest),
		a.AtMapCoord(o, mx, my, tx-1, ty).setHeading(West),
		a.AtMapCoord(o, mx, my, tx-1, ty-1).setHeading(NorthWest),
	}
}

// withinRadius returns all tiles with `r` of (tx, ty) in map (mx, my) that
// matches some func
func (a *Autotiler) withinRadius(o Outline, mx, my, tx, ty, r int, fn func(*Area) bool) []*Area {
	found := []*Area{}

	for iy := ty - r; iy <= ty+r; iy++ {
		for ix := tx - r; ix <= tx+r; ix++ {
			candidate := a.AtMapCoord(o, mx, my, ix, iy)
			if !fn(candidate) {
				continue
			}
			found = append(found, candidate)
		}
	}

	return found
}
