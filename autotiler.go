package autotile

import (
	"image"
	"math/rand"

	"github.com/voidshard/tile"
)

// fnPlacements is used internal to reference the various "placeX" functions
type fnPlacements func(Outline, *rand.Rand, *area, bool) ([]*Event, string, error)

// Autotiler is a struct that understands how to place various sets of tiles in in order
// to create a tiled map.
// - Only "orthogonal" square tiled maps are supported.
type Autotiler struct {
	cfg *Config

	evt chan *Event
	out chan *Event
}

// NewAutotiler creates & returns an autotiler object.
func NewAutotiler(cfg *Config) (*Autotiler, error) {
	at := &Autotiler{cfg: cfg, evt: make(chan *Event)}

	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	go func() { // event pump
		for e := range at.evt {
			if at.out == nil {
				// drop event if no one is listening
				continue
			}
			at.out <- e
		}
	}()

	return at, nil
}

// emitEvent pushes to our internal event chan
func (a *Autotiler) emitEvent(e *Event) {
	a.evt <- e
}

// Events allows a user to see what the autotiler is placing where
// as the decisions are made.
// If Events() is called the caller must consume events.
func (a *Autotiler) Events() <-chan *Event {
	if a.out == nil {
		a.out = make(chan *Event)
	}
	return a.out
}

// SetLand uses data from the Outline `o` to place tiles for the region `region` on the
// given Tileable map `t`
// This sets
// - grass, sand, dirt, snow, rock tiles
// - water, lava
// - cliffs, waterfalls
func (a *Autotiler) SetLand(o Outline, region image.Rectangle, t tile.Tileable) error {
	enact := func(events []*Event) error {
		if events == nil {
			return nil
		}
		for _, e := range events {
			err := t.Set(e.X, e.Y, e.Z, e.Src)
			if err != nil {
				return err
			}
			err = t.SetProperties(e.Src, e.Properties)
			if err != nil {
				return err
			}
			a.emitEvent(e)
		}
		return nil
	}

	placementFuncs := []fnPlacements{a.placeNull, a.placeLand, a.placeWater, a.placeRoad, a.placeMolten, a.placeCliffs}

	seed := a.cfg.Seed + int64(region.Min.X)*int64(region.Max.Y) - int64(region.Max.X)*int64(region.Min.Y)
	rng := rand.New(rand.NewSource(seed))
	var err error
	var evts []*Event

	for ty := region.Min.Y; ty < region.Max.Y; ty++ {
		for tx := region.Min.X; tx < region.Max.X; tx++ {
			me := newArea(o, tx, ty)
			for _, pfn := range placementFuncs {
				evts, _, err = pfn(o, rng, me, false)
				if err != nil {
					return err
				}
				err = enact(evts)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (a *Autotiler) placeNull(o Outline, rng *rand.Rand, me *area, tagsonly bool) ([]*Event, string, error) {
	if !me.Data.IsNull() {
		return nil, "", nil
	}
	if tagsonly {
		return nil, Null, nil
	}

	tiles := me.Data.Tiles()
	if tiles == nil || tiles.Null == "" {
		return nil, Null, nil
	}

	return []*Event{newEvent(me.X, me.Y, a.cfg.ZOffsetLand, tiles.Null, propertiesNull)}, Null, nil
}

func (a *Autotiler) placeLand(o Outline, rng *rand.Rand, me *area, _ bool) ([]*Event, string, error) {
	tiles := me.Data.Tiles()
	if tiles == nil {
		return nil, "", nil
	}

	src := firstFull(rng, tiles.Grass, tiles.Dirt, tiles.Rock)
	tag := Grass
	srcT := ""

	beach := a.cfg.BeachWidth
	nearWtr := false
	nearWtrPlus := false
	if beach > 0 {
		nearWtr = len(withinRadius(o, me.X, me.Y, beach, func(in *area) bool { return in.Data.IsWater() })) > 0
		nearWtrPlus = len(withinRadius(o, me.X, me.Y, beach+1, func(in *area) bool { return in.Data.IsWater() })) > 0
	}
	tsn := a.cfg.TransitionWidth

	if me.Data.Temperature() <= a.cfg.SnowLevel-tsn {
		src = firstFull(rng, tiles.Snow, tiles.Dirt, tiles.Rock)
		tag = Snow
	} else if me.Data.Temperature() <= a.cfg.SnowLevel {
		srcT = firstTransition(rng, tiles.Snow, tiles.Dirt, tiles.Rock)
	} else if nearWtrPlus && !nearWtr {
		srcT = firstTransition(rng, tiles.Sand, tiles.Rock)
	} else if nearWtr {
		src = firstFull(rng, tiles.Sand, tiles.Rock)
		tag = Sand
	} else if me.Data.Height() >= a.cfg.MountainLevel+5 {
		src = firstFull(rng, tiles.Rock, tiles.Dirt)
		tag = Rock
	} else if me.Data.Height() >= a.cfg.MountainLevel {
		srcT = firstTransition(rng, tiles.Rock, tiles.Dirt)
	} else if me.Data.Temperature() <= a.cfg.VegetationMinTemp-tsn {
		src = firstFull(rng, tiles.Dirt, tiles.Rock)
		tag = Dirt
	} else if me.Data.Temperature() <= a.cfg.VegetationMinTemp {
		srcT = firstTransition(rng, tiles.Dirt, tiles.Rock)
	} else if me.Data.Temperature() >= a.cfg.VegetationMaxTemp+tsn { // desert
		src = firstFull(rng, tiles.Sand, tiles.Rock, tiles.Dirt)
		tag = Sand
	} else if me.Data.Temperature() >= a.cfg.VegetationMaxTemp {
		srcT = firstTransition(rng, tiles.Sand, tiles.Rock, tiles.Dirt)
	}

	ret := []*Event{newEvent(me.X, me.Y, a.cfg.ZOffsetLand, src, propertiesLand)}
	if srcT != "" {
		ret = append(ret, newEvent(me.X, me.Y, a.cfg.ZOffsetLand+1, srcT, propertiesLand))
	}

	return ret, tag, nil
}

func (a *Autotiler) placeWater(o Outline, rng *rand.Rand, me *area, tagonly bool) ([]*Event, string, error) {
	if !me.Data.IsWater() {
		return nil, "", nil
	}

	tiles := me.Data.Tiles()
	if tiles == nil || tiles.Water == nil {
		return nil, "", nil
	}
	if tagonly {
		return nil, Water, nil
	}

	crd := cardinals(o, me.X, me.Y)
	in := []*area{}
	out := []*area{}
	for _, t := range crd.all() {
		if t.Data.IsWater() {
			in = append(in, t)
		} else {
			out = append(out, t)
		}
	}

	src := tiles.Water.choosePiece(rng, in, out)
	return []*Event{
		newEvent(me.X, me.Y, a.cfg.ZOffsetWater, src, propertiesWater),
	}, Water, nil
}

func (a *Autotiler) placeRoad(o Outline, rng *rand.Rand, me *area, tagonly bool) ([]*Event, string, error) {
	if !me.Data.IsLand() || !me.Data.IsRoad() {
		return nil, "", nil
	}
	tiles := me.Data.Tiles()
	if tiles.Road == nil {
		return nil, "", nil
	}
	if tagonly {
		return nil, Road, nil
	}

	crd := cardinals(o, me.X, me.Y)
	in := []*area{}
	out := []*area{}
	for _, t := range crd.all() {
		if t.Data.IsRoad() && t.Data.IsLand() {
			in = append(in, t)
		} else {
			out = append(out, t)
		}
	}

	src := tiles.Road.choosePiece(rng, in, out)
	return []*Event{
		newEvent(me.X, me.Y, a.cfg.ZOffsetRoad, src, propertiesRoad),
	}, Road, nil
}

func (a *Autotiler) placeMolten(o Outline, rng *rand.Rand, me *area, tagonly bool) ([]*Event, string, error) {
	if !me.Data.IsMolten() || me.Data.IsWater() {
		return nil, "", nil
	}
	tiles := me.Data.Tiles()
	if tiles.Lava == nil {
		return nil, "", nil
	}
	if tagonly {
		return nil, Lava, nil
	}

	molten := []*area{}
	notmolten := []*area{}
	crd := cardinals(o, me.X, me.Y)
	for _, t := range crd.all() {
		if t.Data.IsMolten() && !t.Data.IsWater() {
			molten = append(molten, t)
		} else {
			notmolten = append(notmolten, t)
		}
	}

	src := tiles.Lava.choosePiece(rng, molten, notmolten)
	return []*Event{
		newEvent(me.X, me.Y, a.cfg.ZOffsetWater, src, propertiesLava),
	}, Lava, nil

}

func (a *Autotiler) placeCliffs(o Outline, rng *rand.Rand, me *area, tagonly bool) ([]*Event, string, error) {
	if me.Data.Height() < a.cfg.CliffLevel {
		return nil, "", nil
	}
	tiles := me.Data.Tiles()
	if tiles.Cliff == nil {
		return nil, "", nil
	}

	crd := cardinals(o, me.X, me.Y)

	lowland := crd.Lower(me.Data.Height())
	if len(lowland) == 0 {
		belowTags, err := a.TagsAt(o, me.X, me.Y+1)
		if err != nil {
			return nil, "", err
		}
		if contains(CliffFace, belowTags) {
			return nil, CliffEdge, nil
		}

		return nil, "", nil
	}
	if tagonly {
		return nil, CliffFace, nil
	}

	c := tiles.Cliff
	evts := []*Event{}
	low := headings(lowland)
	if len(lowland) == 1 {
		switch lowland[0].heading {
		case SouthEast:
			evts = append(evts, newEvent(me.X, me.Y-1, a.cfg.ZOffsetCliff, one(rng, c.ThreeQuarterNorthWest), propertiesCliff))
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetCliff, one(rng, c.WestHalf), propertiesCliff))
		case SouthWest:
			evts = append(evts, newEvent(me.X, me.Y-1, a.cfg.ZOffsetCliff, one(rng, c.ThreeQuarterNorthEast), propertiesCliff))
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetCliff, one(rng, c.EastHalf), propertiesCliff))
		}
	} else if len(lowland) >= 2 || len(lowland) <= 5 {
		if includes(low, cornerNE...) {
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetCliff, one(rng, c.QuarterSouthWest), propertiesCliff))
		} else if includes(low, cornerNW...) {
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetCliff, one(rng, c.QuarterSouthEast), propertiesCliff))
		} else if includes(low, cornerSE...) {
			evts = append(evts, newEvent(me.X, me.Y-1, a.cfg.ZOffsetCliff, one(rng, c.QuarterNorthWest), propertiesCliff))
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetCliff, one(rng, c.QuarterNorthWestBase), propertiesCliff))
		} else if includes(low, cornerSW...) {
			evts = append(evts, newEvent(me.X, me.Y-1, a.cfg.ZOffsetCliff, one(rng, c.QuarterNorthEast), propertiesCliff))
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetCliff, one(rng, c.QuarterNorthEastBase), propertiesCliff))
		} else if includes(low, North) {
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetCliff, one(rng, c.SouthHalf), propertiesCliff))
		} else if includes(low, East) {
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetCliff, one(rng, c.WestHalf), propertiesCliff))
		} else if includes(low, South) {
			evts = append(evts, newEvent(me.X, me.Y-1, a.cfg.ZOffsetCliff, one(rng, c.NorthHalf), propertiesCliff))
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetCliff, one(rng, c.NorthHalfBase), propertiesCliff))
		} else if includes(low, West) {
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetCliff, one(rng, c.EastHalf), propertiesCliff))
		}
	}

	// now we deal with waterfalls ..
	if !me.Data.IsWater() || tiles.Waterfall == nil || len(lowland) <= 1 {
		return evts, CliffFace, nil
	}
	if includes(low, cornerNE...) || includes(low, cornerNW...) || includes(low, cornerSE...) || includes(low, cornerSW...) {
		// we don't place waterfalls on cliff corners
		return evts, CliffFace, nil
	}
	watertiles := []*area{}
	for _, t := range crd.all() {
		if t.Data.IsWater() {
			watertiles = append(watertiles, t)
		}
	}
	if len(watertiles) < 5 {
		return evts, CliffFace, nil
	}

	if includes(low, North) { // waterfall flowing north
		if tiles.Waterfall.SN == nil {
			return evts, CliffFace, nil
		}
		if len(watertiles) < 8 {
			return evts, CliffFace, nil
		}
		// we can only see the top of the waterfall falling out of view
		evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetWaterfall, one(rng, tiles.Waterfall.SN.MidTop), propertiesWFall))
	} else if includes(low, South) { // waterfall flowing south
		if tiles.Waterfall.NS == nil {
			return evts, CliffFace, nil
		}

		// check the tile above (Y-1) of this tile
		above := newArea(o, me.X, me.Y-1)
		abovelowland := cardinals(o, me.X, me.Y-1).Lower(above.Data.Height())
		istop := len(abovelowland) < 2

		wet := headings(watertiles)
		if len(watertiles) >= 8 { // middle
			if istop {
				evts = append(evts, newEvent(me.X, me.Y-1, a.cfg.ZOffsetWaterfall, one(rng, tiles.Waterfall.NS.MidTop), propertiesWFall))
			}
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetWaterfall, one(rng, tiles.Waterfall.NS.MidCentre), propertiesWFall))
			evts = append(evts, newEvent(me.X, me.Y+1, a.cfg.ZOffsetWaterfall, one(rng, tiles.Waterfall.NS.MidBottom), propertiesWFall))
		} else if includes(wet, cornerNE...) && includes(wet, cornerSE...) { // left
			if istop {
				evts = append(evts, newEvent(me.X, me.Y-1, a.cfg.ZOffsetWaterfall, one(rng, tiles.Waterfall.NS.LeftTop), propertiesWFall))
			}
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetWaterfall, one(rng, tiles.Waterfall.NS.LeftCentre), propertiesWFall))
			evts = append(evts, newEvent(me.X, me.Y+1, a.cfg.ZOffsetWaterfall, one(rng, tiles.Waterfall.NS.LeftBottom), propertiesWFall))
		} else if includes(wet, cornerNW...) && includes(wet, cornerSW...) { // right
			if istop {
				evts = append(evts, newEvent(me.X, me.Y-1, a.cfg.ZOffsetWaterfall, one(rng, tiles.Waterfall.NS.RightTop), propertiesWFall))
			}
			evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetWaterfall, one(rng, tiles.Waterfall.NS.RightCentre), propertiesWFall))
			evts = append(evts, newEvent(me.X, me.Y+1, a.cfg.ZOffsetWaterfall, one(rng, tiles.Waterfall.NS.RightBottom), propertiesWFall))
		}
	}

	return evts, CliffFace, nil
}

// TagsAt indicates the 'tags' that will be set a given location, roughly
// indicating the kind of tile(s) that will be placed there.
// We return user set tags + the tag we consider the most important for the terrain.
func (a *Autotiler) TagsAt(o Outline, x, y int) ([]string, error) {
	me := newArea(o, x, y)
	placementFuncs := []fnPlacements{a.placeNull, a.placeCliffs, a.placeWater, a.placeMolten, a.placeRoad, a.placeLand}
	rng := rand.New(rand.NewSource(0)) // doesn't affect the tag we get (only the specific src image)

	var err error
	var tag string
	for _, pfn := range placementFuncs {
		// passing true here to skip logic that doesn't change the final tag
		_, tag, err = pfn(o, rng, me, true)
		if err != nil {
			return nil, err
		}
		if tag != "" {
			break // accept the first valid tag
		}
	}

	tags := me.Data.Tags()
	if tags == nil {
		tags = []string{tag} // nb. tag can be Null
	} else {
		tags = append(tags, tag)
	}

	return tags, nil
}

// SetObjects places objects from the given ObjectBin on to the map `t` within the area defined by
// the region.
func (a *Autotiler) SetObjects(o Outline, region image.Rectangle, t tile.Tileable, bin ObjectBin) error {
	for ty := region.Min.Y; ty < region.Max.Y; ty++ {
		for tx := region.Min.X; tx < region.Max.X; tx++ {
			// choose an object
			id, obj, err := bin.Choose(t, tx, ty, a.cfg.ZOffsetObject)
			if err != nil {
				return err
			}
			if obj == nil {
				continue
			}

			// place an object!
			err = t.Add(tx, ty, a.cfg.ZOffsetObject, obj)
			if err != nil {
				return err
			}

			// and make an event
			a.emitEvent(newObjEvent(tx, ty, a.cfg.ZOffsetObject, id))
		}
	}
	return nil
}
