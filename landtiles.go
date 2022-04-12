package autotile

import (
	"github.com/voidshard/tile"

	"image"
	"math/rand"
	"sort"
)

// Land comprises tiles that make up the base layers
// of our maps. Grass, water, dirt, rock, sand etc.
//
// Note that all of these contain image paths that are placed in our .tmx
// map & tileset
type LandTiles struct {
	// Null tile is set (or not used if "") where nothing else appliess
	// Eg. in an interior it might be a tile representing
	// "not a valid area" or .. whatever else.
	Null string

	// Grass is a general ground tile, set by default if nothing else applies
	Grass *Tileset

	// Sand is used for deserts, beaches etc
	Sand *Tileset

	// Dirt is placed as a fallback, or when grass cannot be placed & nothing
	// else applies
	Dirt *Tileset

	// Snow is placed instead of grass when temperature is low
	Snow *Tileset

	// Rock is placed generally instead of dirt, or when we're high up (barren mountains)
	Rock *Tileset

	// Water is placed where ever there is .. well .. water
	Water *Tileset

	// Bridge composed of tiles (stepping stones, planks or something that tiles well
	// is recommended). Placed where water + road intersect.
	Bridge *Tileset

	// Road is placed where road is true
	Road *Tileset

	// Lava is placed where molten is true
	Lava *Tileset

	// Cliff is placed above config.CliffLevel where the land changes height.
	// Note we only place cliffs every few height changes.
	// TODO: improve so we place cliffs only with height changes above a
	// certain sharpness.
	Cliff *Tileset

	// Stairs are placed where roads meet cliffs.
	StairsWestEast *Tileset

	//
	StairsEastWest *Tileset

	//
	StairsNorthSouth *Tileset

	//
	StairsSouthNorth *Tileset

	// Waterfall is where a cliff & a river intersect one another.
	WaterfallNorthSouth *Tileset

	// Waterfall flowing up from the bottom (away from us)
	WaterfallSouthNorth *Tileset

	// Waterfall flowing East - West (right -> left)
	WaterfallEastWest *Tileset

	// Waterfall flowing West - East (left -> right)
	WaterfallWestEast *Tileset
}

// Validate all land tile(s) are set
func (b *LandTiles) Validate() error {
	return nil
}

// Tileset represents tiles that must be placed depending on
// what tile(s) of the same terrain are adjacent to them.
// Ie. water or road tiles that are supposed to fit together to form
// a long river or road.
// These are expected to be paths to specific images not map files.
type Tileset struct {
	// Full represents a full tile of our complex type.
	// Ie. a tile covered entirely in water
	Full []string

	// Represents a tile transitioning to this tileset
	Transition []string

	// The North half the tile is our complex type.
	// Ie. the top / North side of tile is water.
	NorthHalf []string

	// The East half the tile is our complex type.
	EastHalf []string

	// The Sough half the tile is our complex type.
	SouthHalf []string

	// The West half the tile is our complex type.
	WestHalf []string

	// 1/4 of the tile is the complex type, that being the NE corner.
	QuarterNorthEast []string

	// 1/4 of the tile is the complex type, that being the SE corner.
	QuarterSouthEast []string

	// 1/4 of the tile is the complex type, that being the SW corner.
	QuarterSouthWest []string

	// 1/4 of the tile is the complex type, that being the NW corner.
	QuarterNorthWest []string

	// 3/4 of the tile is the complex type, centred on the NE corner.
	// Ie, if this type is 'water' then the SouthWest corner here is
	// *not* water (since it's the 1/4 that is *not* our type)
	ThreeQuarterNorthEast []string

	// 3/4 of the tile is the complex type, centred on the SE corner.
	ThreeQuarterSouthEast []string

	// 3/4 of the tile is the complex type, centred on the SW corner.
	ThreeQuarterSouthWest []string

	// 3/4 of the tile is the complex type, centred on the NW corner.
	ThreeQuarterNorthWest []string
}

//
func (t *Tileset) fillRect(rng *rand.Rand, in image.Rectangle, z int, props *tile.Properties) []*Event {
	evts := []*Event{}

	for y := in.Max.Y; y >= in.Min.Y; y-- {
		for x := in.Min.X; x <= in.Max.X; x++ {
			src := ""
			switch y {
			case in.Max.Y: // place max y first so we tile from obj bottom
				switch x {
				case in.Min.X:
					src = one(rng, t.QuarterNorthEast)
				case in.Max.X:
					src = one(rng, t.QuarterNorthWest)
				default:
					src = one(rng, t.NorthHalf)
				}
			case in.Min.Y:
				switch x {
				case in.Min.X:
					src = one(rng, t.QuarterSouthEast)
				case in.Max.X:
					src = one(rng, t.QuarterSouthWest)
				default:
					src = one(rng, t.SouthHalf)
				}
			default:
				switch x {
				case in.Min.X:
					src = one(rng, t.EastHalf)
				case in.Max.X:
					src = one(rng, t.WestHalf)
				default:
					src = one(rng, t.Full)
				}
			}
			if src != "" {
				evts = append(evts, newEvent(x, y, z, src, props))
			}
		}
	}

	return evts
}

// choosePiece decides which piece of our ComplexLand to place given which of the
// surrounding pieces are 'isSet' (this complex land) or 'notSet' (not this).
func (t *Tileset) choosePiece(rng *rand.Rand, crd *nearby, isIn func(a *area) bool) string {
	isSet := []*area{}
	notSet := []*area{}

	in := []Heading{}
	out := []Heading{}

	for _, a := range crd.all() {
		if isIn(a) {
			isSet = append(isSet, a)
			in = append(in, a.heading)
		} else {
			notSet = append(notSet, a)
			out = append(out, a.heading)
		}
	}

	sort.Slice(in, func(i, j int) bool { return int(in[i]) < int(in[j]) })
	sort.Slice(out, func(i, j int) bool { return int(out[i]) < int(out[j]) })

	if len(isSet) >= 8 { // we can't actually have more than 8
		return one(rng, t.Full)
	} else if len(isSet) == 7 { // implies notSet len == 1
		lt := notSet[0] // ie. this is the only element
		switch lt.heading {
		case NorthEast:
			return one(rng, t.ThreeQuarterSouthWest)
		case SouthEast:
			return one(rng, t.ThreeQuarterNorthWest)
		case SouthWest:
			return one(rng, t.ThreeQuarterNorthEast)
		case NorthWest:
			return one(rng, t.ThreeQuarterSouthEast)
		case North:
			return one(rng, t.SouthHalf)
		case East:
			return one(rng, t.WestHalf)
		case South:
			return one(rng, t.NorthHalf)
		case West:
			return one(rng, t.EastHalf)
		}
	} else if len(isSet) <= 6 {
		// try to place corners first, then edges ..
		if includes(out, cornerNE...) {
			if includes(in, edgeSW...) {
				return one(rng, t.ThreeQuarterSouthWest)
			} else {
				return one(rng, t.QuarterSouthWest)
			}
		} else if includes(out, cornerSE...) {
			if includes(in, edgeNW...) {
				return one(rng, t.ThreeQuarterNorthWest)
			} else {
				return one(rng, t.QuarterNorthWest)
			}
		} else if includes(out, cornerSW...) {
			if includes(in, edgeNE...) {
				return one(rng, t.ThreeQuarterNorthEast)
			} else {
				return one(rng, t.QuarterNorthEast)
			}
		} else if includes(out, cornerNW...) {
			if includes(in, edgeSE...) {
				return one(rng, t.ThreeQuarterSouthEast)
			} else {
				return one(rng, t.QuarterSouthEast)
			}
		} else if includes(out, North) {
			return one(rng, t.SouthHalf)
		} else if includes(out, East) {
			return one(rng, t.WestHalf)
		} else if includes(out, South) {
			return one(rng, t.NorthHalf)
		} else if includes(out, West) {
			return one(rng, t.EastHalf)
		}
	}

	return ""
}
