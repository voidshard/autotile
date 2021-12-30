package autotile

import (
	"fmt"
	"math/rand"
)

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
func (c *ComplexLand) choosePiece(rng *rand.Rand, isSet []*Area, notSet []*Area) string {
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
