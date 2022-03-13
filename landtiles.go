package autotile

import (
	"fmt"

	"math/rand"
)

// Land comprises tiles that make up the base layers
// of our maps. Grass, water, dirt, rock, sand etc.
//
// Note that all of these contain image paths that are placed in our .tmx
// map & tileset
type LandTiles struct {
	// Null tile is set (or not used if "") where nothing else applies.
	// Eg. in an interior it might be a tile representing
	// "not a valid area" or .. whatever else.
	Null string

	// Grass is a general ground tile, set by default if nothing else applies
	Grass *BasicLand

	// Sand is used for deserts, beaches etc
	Sand *BasicLand

	// Dirt is placed as a fallback, or when grass cannot be placed & nothing
	// else applies
	Dirt *BasicLand

	// Snow is placed instead of grass when temperature is low
	Snow *BasicLand

	// Rock is placed generally instead of dirt, or when we're high up (barren mountains)
	Rock *BasicLand

	// Water is placed anywhere one of River, Sea or River is true.
	Water *ComplexLand

	// Road is placed where road is true
	Road *ComplexLand

	// Lava is placed where molten is true
	Lava *ComplexLand

	// Cliff is placed above config.CliffLevel where the land changes height.
	// Note we only place cliffs every few height changes.
	// TODO: improve so we place cliffs only with height changes above a
	// certain sharpness.
	Cliff *Cliff

	// Waterfall is where a cliff & a river intersect one another.
	// Nb. Since we only support square orthogonal maps we only draw waterfalls
	// flowing North-South or South-North (flowing straight towards or straight
	// away from us).
	// We could in theory support more, with more tiles.
	Waterfall *Waterfall
}

// BasicLand represents the lowest layer of the map; tiles that
// are expected to decorate the floor.
type BasicLand struct {
	// Tiles representing the base land; dirt, grass etc
	// Expect path(s) to image file(s)
	Full []string

	// Tiles representing partial covering by this land type
	// eg. we might render a grass Base + a snow Transition in
	// order to depict a transition region from grass to snow.
	// Expect path(s) to image file(s)
	Transition []string
}

// Validate all land tile(s) are set
func (b *LandTiles) Validate() error {
	for k, l := range map[string]*BasicLand{
		Grass: b.Grass,
		Sand:  b.Sand,
		Dirt:  b.Dirt,
		Snow:  b.Snow,
		Rock:  b.Rock,
	} {
		if l == nil {
			continue
		}
		if len(l.Full) < 1 {
			return fmt.Errorf("%w: land %v requires full tiles", ErrMissingRequiredValue, k)
		}
	}

	for _, l := range []*ComplexLand{b.Water, b.Road, b.Lava} {
		if l == nil {
			continue
		}

		err := l.validate()
		if err != nil {
			return err
		}
	}

	if b.Cliff != nil {
		err := b.Cliff.validate()
		if err != nil {
			return err
		}
	}

	if b.Waterfall != nil {
		err := b.Waterfall.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// ComplexLand represents tiles that must be placed depending on
// what tile(s) of the same terrain are adjacent to them.
// Ie. water or road tiles that are supposed to fit together to form
// a long river or road.
// These are expected to be paths to specific images not map files.
//
// Water, road & lava all qualify as complex land. Cliffs & waterfalls
// are sufficiently different to have their own structs.
type ComplexLand struct {
	// Full represents a full tile of our complex type.
	// Ie. a tile covered entirely in water
	Full []string

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

// choosePiece decides which piece of our ComplexLand to place given which of the
// surrounding pieces are 'isSet' (this complex land) or 'notSet' (not this).
func (c *ComplexLand) choosePiece(rng *rand.Rand, isSet []*area, notSet []*area) string {
	// Essentially the premise here is around full, 1/2, 1/4 and 3/4 tiles and
	// what to place at the location of 'me' based on what is around us (ie,
	// whether the tiles around us are the same as 'me' -- part of the ComplexLand
	// -- or something else).
	// That, is if this tile is part of ComplexLand (ie "water") and the tiles on
	// the East are not (ie "land") then we want to place our "WestHalf" --
	// because we're saying that that WestHalf of our tile should have our ComplexLand
	// (in this case "water").

	in := headings(isSet)
	out := headings(notSet)

	src := one(rng, c.Full)

	if len(isSet) >= 8 { // we can't actually have more than 8
		return src
	} else if len(isSet) == 7 { // implies notSet len == 1
		lt := notSet[0]
		switch lt.heading {
		case NorthEast:
			src = one(rng, c.ThreeQuarterSouthWest)
		case SouthEast:
			src = one(rng, c.ThreeQuarterNorthWest)
		case SouthWest:
			src = one(rng, c.ThreeQuarterNorthEast)
		case NorthWest:
			src = one(rng, c.ThreeQuarterSouthEast)
		case North:
			src = one(rng, c.SouthHalf)
		case East:
			src = one(rng, c.WestHalf)
		case South:
			src = one(rng, c.NorthHalf)
		case West:
			src = one(rng, c.EastHalf)
		}
	} else if len(isSet) <= 6 {
		// try to place corners first, then edges ..
		if includes(out, cornerNE...) {
			if includes(in, edgeSW...) {
				src = one(rng, c.ThreeQuarterSouthWest)
			} else {
				src = one(rng, c.QuarterSouthWest)
			}
		} else if includes(out, cornerSE...) {
			if includes(in, edgeNW...) {
				src = one(rng, c.ThreeQuarterNorthWest)
			} else {
				src = one(rng, c.QuarterNorthWest)
			}
		} else if includes(out, cornerSW...) {
			if includes(in, edgeNE...) {
				src = one(rng, c.ThreeQuarterNorthEast)
			} else {
				src = one(rng, c.QuarterNorthEast)
			}
		} else if includes(out, cornerNW...) {
			if includes(in, edgeSE...) {
				src = one(rng, c.ThreeQuarterSouthEast)
			} else {
				src = one(rng, c.QuarterSouthEast)
			}
		} else if includes(out, North) {
			src = one(rng, c.SouthHalf)
		} else if includes(out, East) {
			src = one(rng, c.WestHalf)
		} else if includes(out, South) {
			src = one(rng, c.NorthHalf)
		} else if includes(out, West) {
			src = one(rng, c.EastHalf)
		}
	}

	return src
}

// validate all fields are set
func (c *ComplexLand) validate() error {
	for k, v := range map[string][]string{
		"full":           c.Full,
		"north-half":     c.NorthHalf,
		"east-half":      c.EastHalf,
		"south-half":     c.SouthHalf,
		"west-half":      c.WestHalf,
		"1/4-north-east": c.QuarterNorthEast,
		"1/4-south-east": c.QuarterSouthEast,
		"1/4-south-west": c.QuarterSouthWest,
		"1/4-north-west": c.QuarterNorthWest,
		"3/4-north-east": c.ThreeQuarterNorthEast,
		"3/4-south-east": c.ThreeQuarterSouthEast,
		"3/4-south-west": c.ThreeQuarterSouthWest,
		"3/4-north-west": c.ThreeQuarterNorthWest,
	} {
		if v == nil || len(v) < 1 {
			return fmt.Errorf("%w: %s requires at least one tile", ErrMissingRequiredValue, k)
		}
	}

	return nil
}

// Cliff is a variant of ComplexLand with a twist.
// Note here that we're refering here always to the parts of the tile
// that are low ground.
// Ie. the 'NorthHalf' tile means that the *SouthHalf* of the tile
// is raised land (we mean "this image should form the North part
// of the finished cliff").
// The Southern pieces all include "Base" pieces that go one tile
// beneath (y+1) their counterpart(s) to make the cliff appear high.
type Cliff struct {
	// NorthHalf implies that the north half of the tile is high ground.
	// Ie. the cliff should be rising up from South -> North.
	NorthHalf []string

	// NorthHalfBase is the bottom piece of the NorthHalf cliff
	NorthHalfBase []string

	// EastHalf implies the land is rising West -> East.
	EastHalf []string

	// SouthHalf is a cliff facing away from the user, falling off to the North.
	SouthHalf []string

	// WestHalf is a cliff rising East -> West
	WestHalf []string

	// QuarterNorthEast implies that the North, East & North-East parts of the tile
	// are high ground (ie. the cliff is turning a corner).
	QuarterNorthEast []string

	// QuarterNorthEastBase is the lower corner half of QuarterNorthEast
	QuarterNorthEastBase []string

	// QuarterSouthEast implies South, East & South-East are high ground (cliff corner)
	QuarterSouthEast []string

	// QuarterSouthWest implies South, West & South-West are high ground (cliff corner)
	QuarterSouthWest []string

	// QuarterNorthWest implies that the North, West & North-West parts of the tile
	// are high ground (ie. the cliff is turning a corner).
	QuarterNorthWest []string

	// QuarterNorthWestBase is the bottom half of a QuarterNorthWest cliff tile.
	QuarterNorthWestBase []string

	// ThreeQuarterNorthEast implies that the *low* sections are South, West and South-West
	// ie, the South-West corner is the bottom of a cliff.
	ThreeQuarterNorthEast []string

	// ThreeQuarterNorthWest implies that the *low* sections are South, East and South-East
	// ie, the South-East corner is the bottom of a cliff.
	ThreeQuarterNorthWest []string
}

func (c *Cliff) validate() error {
	for k, v := range map[string][]string{
		"north-half":          c.NorthHalf,
		"north-half-base":     c.NorthHalfBase,
		"east-half":           c.EastHalf,
		"south-half":          c.SouthHalf,
		"west-half":           c.WestHalf,
		"1/4-north-east":      c.QuarterNorthEast,
		"1/4-north-east-base": c.QuarterNorthEastBase,
		"1/4-south-east":      c.QuarterSouthEast,
		"1/4-south-west":      c.QuarterSouthWest,
		"1/4-north-west":      c.QuarterNorthWest,
		"1/4-north-west-base": c.QuarterNorthWestBase,
		"3/4-north-east":      c.ThreeQuarterNorthEast,
		"3/4-north-west":      c.ThreeQuarterNorthWest,
	} {
		if v == nil || len(v) < 1 {
			return fmt.Errorf("%w: %s requires at least one tile", ErrMissingRequiredValue, k)
		}
	}

	return nil
}

// Waterfall is where a river meets a cliff.
// Nb. we only draw waterfalls flowing North -> South or South -> North.
type Waterfall struct {
	// NS -> flowing North-South towards us
	NS *WaterfallNorthSouth

	// SN -> flowing South-North away from us
	SN *WaterfallSouthNorth
}

func (w *Waterfall) validate() error { return nil }

// WaterfallNorthSouth are waterfall tiles for a waterfall flowing South -> North.
// Nb. South -> North implies the waterfall is flowing away from us over a cliff.
type WaterfallSouthNorth struct {
	MidTop []string
}

// WaterfallNorthSouth are waterfall tiles for a waterfall flowing North -> South
// Nb. North -> South implies the waterfall is flowing towards us down a cliff.
type WaterfallNorthSouth struct {
	LeftTop     []string
	LeftCentre  []string
	LeftBottom  []string
	MidTop      []string
	MidCentre   []string
	MidBottom   []string
	RightTop    []string
	RightCentre []string
	RightBottom []string
}
