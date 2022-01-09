package autotile

import (
	"fmt"
	"math/rand"
)

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

// placement is an internal struct so we can return a set of tiles to be placed
type placement struct {
	X   int
	Y   int
	Src string
}

// placements returns cliff tiles we want to set
func (c *Cliff) placements(rng *rand.Rand, me *Area, lowland []*Area) []*placement {
	if len(lowland) == 0 {
		return []*placement{}
	}

	low := headings(lowland)
	if len(lowland) == 1 {
		switch lowland[0].heading {
		case SouthEast:
			return []*placement{
				&placement{X: me.X, Y: me.Y - 1, Src: one(rng, c.ThreeQuarterNorthWest)},
				&placement{X: me.X, Y: me.Y, Src: one(rng, c.WestHalf)},
			}
		case SouthWest:
			return []*placement{
				&placement{X: me.X, Y: me.Y - 1, Src: one(rng, c.ThreeQuarterNorthEast)},
				&placement{X: me.X, Y: me.Y, Src: one(rng, c.EastHalf)},
			}
		}
	} else if len(lowland) >= 2 || len(lowland) <= 5 {
		if includes(low, cornerNE...) {
			return []*placement{
				&placement{X: me.X, Y: me.Y, Src: one(rng, c.QuarterSouthWest)},
			}
		} else if includes(low, cornerNW...) {
			return []*placement{
				&placement{X: me.X, Y: me.Y, Src: one(rng, c.QuarterSouthEast)},
			}
		} else if includes(low, cornerSE...) {
			return []*placement{
				&placement{X: me.X, Y: me.Y - 1, Src: one(rng, c.QuarterNorthWest)},
				&placement{X: me.X, Y: me.Y, Src: one(rng, c.QuarterNorthWestBase)},
			}
		} else if includes(low, cornerSW...) {
			return []*placement{
				&placement{X: me.X, Y: me.Y - 1, Src: one(rng, c.QuarterNorthEast)},
				&placement{X: me.X, Y: me.Y, Src: one(rng, c.QuarterNorthEastBase)},
			}
		} else if includes(low, North) {
			return []*placement{
				&placement{X: me.X, Y: me.Y, Src: one(rng, c.SouthHalf)},
			}
		} else if includes(low, East) {
			return []*placement{
				&placement{X: me.X, Y: me.Y, Src: one(rng, c.WestHalf)},
			}
		} else if includes(low, South) {
			return []*placement{
				&placement{X: me.X, Y: me.Y - 1, Src: one(rng, c.NorthHalf)},
				&placement{X: me.X, Y: me.Y, Src: one(rng, c.NorthHalfBase)},
			}
		} else if includes(low, West) {
			return []*placement{
				&placement{X: me.X, Y: me.Y, Src: one(rng, c.EastHalf)},
			}
		}
	}

	return []*placement{}
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
