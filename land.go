package autotile

import (
	"fmt"
)

// Land comprises tiles that make up the base layers
// of our maps. Grass, water, dirt, rock, sand etc.
//
// Note that all of these contain image paths that are placed in our .tmx
// map & tileset
type Land struct {
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

	// Lava is placed where lava is true
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

	// tag indicates what kind of land this is.
	tag string
}

// setTags ensures land tags are set.
func (b *Land) setTags() {
	if b.Grass != nil {
		b.Grass.tag = Grass
	}
	if b.Sand != nil {
		b.Sand.tag = Sand
	}
	if b.Dirt != nil {
		b.Dirt.tag = Dirt
	}
	if b.Snow != nil {
		b.Snow.tag = Snow
	}
	if b.Rock != nil {
		b.Rock.tag = Rock
	}
}

// Validate all land tile(s) are set
func (b *Land) Validate() error {
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
