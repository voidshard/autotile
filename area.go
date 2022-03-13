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

// area represents a single X,Y on the map
// We use this internally so we can encapsulate some metadata about where this
// came from & it's relation to other points of interest.
type area struct {
	// the X co-ord of this area
	X int

	// the Y co-ord of this area
	Y int

	Data LandData

	// internally used to track this tiles' relation to another
	heading Heading
}

// setHeading sets heading to `in` and returns the area
func (a *area) setHeading(in Heading) *area {
	a.heading = in
	return a
}

// headings returns the `heading` values of the given areas
func headings(in []*area) []Heading {
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

// nearby is a set of areas near (adjacent to) a given tile
type nearby struct {
	North     *area
	NorthEast *area
	East      *area
	SouthEast *area
	South     *area
	SouthWest *area
	West      *area
	NorthWest *area
}

// all returns all tiles nearby in 8 cardinal directions
func (n *nearby) all() []*area {
	return []*area{
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
func (n *nearby) Lower(h int) []*area {
	land := []*area{}
	for _, t := range n.all() {
		if t.Data.Height() < h {
			land = append(land, t)
		}
	}
	return land
}

// cardinals returns information on surrounding tiles
func cardinals(o Outline, tx, ty int) *nearby {
	return &nearby{
		newArea(o, tx, ty-1).setHeading(North),
		newArea(o, tx+1, ty-1).setHeading(NorthEast),
		newArea(o, tx+1, ty).setHeading(East),
		newArea(o, tx+1, ty+1).setHeading(SouthEast),
		newArea(o, tx, ty+1).setHeading(South),
		newArea(o, tx-1, ty+1).setHeading(SouthWest),
		newArea(o, tx-1, ty).setHeading(West),
		newArea(o, tx-1, ty-1).setHeading(NorthWest),
	}
}

// newArea returns a new area of the given co-ord pair
func newArea(o Outline, x, y int) *area {
	return &area{X: x, Y: y, Data: o.LandAt(x, y)}
}

// withinRadius returns all tiles with `r` of (tx, ty) in map (mx, my) that
// matches some func
func withinRadius(o Outline, tx, ty, r int, fn func(*area) bool) []*area {
	found := []*area{}

	for iy := ty - r; iy <= ty+r; iy++ {
		for ix := tx - r; ix <= tx+r; ix++ {
			if tx == ix && ty == iy {
				continue
			}

			candidate := newArea(o, ix, iy)
			if !fn(candidate) {
				continue
			}

			found = append(found, candidate)
		}
	}

	return found
}
