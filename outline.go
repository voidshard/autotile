package autotile

import (
	"image"
	"math/rand"

	"github.com/voidshard/tile"
)

// Outline is rough guide for building a world from.
// Essentially an interface to a zoomed out world map that describes the world in
// a general sense. All funcs should be thread-safe.
type Outline interface {
	// Bounds returns the size of the world space
	Bounds() image.Rectangle

	// At returns specific data about a given point on the world map.
	// Nb. this function should return rapidly as it will be called a *lot*
	// and delays here will cause our map creation to take a *long* time.
	At(x, y int) *Area
}

// MapOutline represents a single tile map along with various bits of
// helpful metadata attached. These are created internally & returned
// to the user.
type MapOutline struct {
	parent *Autotiler

	// MapX is the world map offset (x value)
	MapX int

	// MapY is the world map offset (y value)
	MapY int

	// MapWidth is the width of this map in tiles
	MapWidth int

	// MapHeight is the height of this map in tiles
	MapHeight int

	// indicates tiles we should fill with water (even if they're
	// not marked as such)
	flood []bool

	// indicates tiles we should consider 'road' in addition
	// to marked tiles
	road []bool

	// indicates tiles we flood with lava, similar to water
	lava []bool

	// metadata about the area so we can avoid running code we don't need to
	numhighlands int // ie. land we might place cliffs on
	numwater     int
	numland      int // ie. not water
	numroad      int
	numlava      int

	// rng with our seed
	seed int64
	rng  *rand.Rand

	// Tilemap is the tile.Map that we're setting tiles on.
	// A user should be free to add objects to the map before using an objectbin
	// (as we do not overwrite objects when placing from a bin).
	// A user will need to write out this map. It's recommended that the map
	// is removed from memory reasonably snapily in order to preserve memory.
	Tilemap *tile.Map

	// tags set on each map index
	tags []map[string]bool
}

// Seed returns our map outlines RNG seed value
func (m *MapOutline) Seed() int64 {
	return m.seed
}

// AllTags returns all tags set on any tiles in the given map & their frequency
func (m *MapOutline) AllTags() map[string]int {
	found := map[string]int{}
	for _, m := range m.tags {
		for k, v := range m {
			if v {
				count, _ := found[k]
				found[k] = count + 1
			}
		}
	}
	return found
}

// hasTag returns two bools, the first indicates if the given tag is present
// on the given tile, the second indicates if the given tile (x,y) is valid.
// Similar to a map returning (value, ok).
func (m *MapOutline) hasTag(x, y int, tag string) (bool, bool) {
	if x < 0 || x >= m.MapWidth || y < 0 || y >= m.MapHeight {
		return false, false
	}
	tags := m.tags[m.MapWidth*y+x]
	if tags == nil {
		return false, true
	}
	v, _ := tags[tag]
	return v, true
}

// Tags returns tags at the given map index
func (m *MapOutline) Tags(x, y int) []string {
	if x < 0 || x >= m.MapWidth || y < 0 || y >= m.MapHeight {
		return []string{}
	}

	result := []string{}

	tags := m.tags[m.MapWidth*y+x]

	if tags != nil {
		for t := range tags {
			result = append(result, t)
		}
	}

	return result
}

// SetTags overwrites tags at (x,y) with the given tags
func (m *MapOutline) SetTags(x, y int, tags ...string) {
	if x < 0 || x >= m.MapWidth || y < 0 || y >= m.MapHeight {
		return
	}

	tm := map[string]bool{}
	for _, t := range tags {
		if t == "" {
			continue
		}
		tm[t] = true
	}

	m.tags[m.MapWidth*y+x] = tm
}

// HasTag returns if the given tile (x, y) has the given tag.
// An out of bounds tile will also return false.
func (m *MapOutline) HasTag(x, y int, tag string) bool {
	v, ok := m.hasTag(x, y, tag)
	return ok && v
}

// HasAnyTags is a logical extension of HasTag where we return true if the given
// tile at (x, y) has any of the given tags.
func (m *MapOutline) HasAnyTags(x, y int, tags []string) bool {
	for _, t := range tags {
		v, ok := m.hasTag(x, y, t)
		if !ok {
			// tile is invalid -> answer is false
			return false
		}
		if v {
			return true
		}
	}
	return false
}

// setLava marks the given index for lava flooding
func (m *MapOutline) setLava(x, y int) {
	if x < 0 || x >= m.MapWidth || y < 0 || y >= m.MapHeight {
		return
	}
	m.lava[m.MapWidth*y+x] = true
}

// IsLava returns if the given tile has been marked for lava flooding
func (m *MapOutline) IsLava(x, y int) bool {
	if x < 0 || x >= m.MapWidth || y < 0 || y >= m.MapHeight {
		return false
	}
	return m.lava[m.MapWidth*y+x]
}

// setRoad sets the given tile for `flooding` with road
func (m *MapOutline) setRoad(x, y int) {
	if x < 0 || x >= m.MapWidth || y < 0 || y >= m.MapHeight {
		return
	}
	m.road[m.MapWidth*y+x] = true
}

// IsRoad returns if the given tile is marked for `flooding`
func (m *MapOutline) IsRoad(x, y int) bool {
	if x < 0 || x >= m.MapWidth || y < 0 || y >= m.MapHeight {
		return false
	}
	return m.road[m.MapWidth*y+x]
}

// setWater will mark this tile as 'water' regardless of the result from Outline
func (m *MapOutline) setWater(x, y int) {
	if x < 0 || x >= m.MapWidth || y < 0 || y >= m.MapHeight {
		return
	}
	m.flood[m.MapWidth*y+x] = true
}

// IsWater returns if we've marked this tile for flooding (ie. as water).
// This happens because we find abrupt diagonal transitions of river / coastline
// problematic with the available tiles, so we adjust things for rendering
// purposes.
func (m *MapOutline) IsWater(x, y int) bool {
	if x < 0 || x >= m.MapWidth || y < 0 || y >= m.MapHeight {
		return false
	}
	return m.flood[m.MapWidth*y+x]
}
