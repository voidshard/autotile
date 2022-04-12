package autotile

import (
	"fmt"
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
	// for some objects that involve intersections of tiles we mark collision
	// squares and come back to them later
	collisions := newCollisionHandler()

	enact := func(events []*Event) error {
		if events == nil {
			return nil // nothing to do
		}

		for _, e := range events {
			if e.collisionType != "" {
				collisions.append(e)
				continue // collisions are internal information
			}
			if e.Src == "" && e.ObjectID == "" {
				continue // nothing is set
			}
			if e.X < region.Min.X || e.X >= region.Max.X || e.Y < region.Min.Y && e.Y >= region.Max.Y {
				continue // outside of the area
			}

			err := t.Set(e.X, e.Y, e.Z, e.Src)
			if err != nil {
				return err
			}
			err = t.SetProperties(e.Src, e.Properties)
			if err != nil {
				return err
			}

			a.emitEvent(e) // now that we've done it, push to listener (if any)
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

	for _, col := range collisions.All() {
		evts = a.handleCollision(o, rng, col)
		err = enact(evts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Autotiler) handleCollision(o Outline, rng *rand.Rand, col *collision) []*Event {
	e := col.Events()[0]
	me := newArea(o, e.X, e.Y)
	tiles := me.Data.Tiles()
	r := col.Max()

	evts := []*Event{}

	fmt.Println("DETECTED", col.typ, r)

	switch col.typ {
	case collisionStairsNS:
		if tiles.StairsNorthSouth == nil {
			return nil
		}
		evts = tiles.StairsNorthSouth.fillRect(
			rng,
			image.Rect(r.Min.X, r.Min.Y-1, r.Max.X, r.Max.Y+1),
			a.cfg.ZOffsetWaterfall,
			propertiesRoad,
		)
	case collisionStairsSN:
		if tiles.StairsSouthNorth == nil {
			return nil
		}
		evts = tiles.StairsSouthNorth.fillRect(rng, r, a.cfg.ZOffsetWaterfall, propertiesRoad)
	case collisionStairsEW:
		if tiles.StairsEastWest == nil {
			return nil
		}
		evts = tiles.StairsEastWest.fillRect(
			rng,
			image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y-1),
			a.cfg.ZOffsetWaterfall,
			propertiesRoad,
		)
	case collisionStairsWE:
		if tiles.StairsWestEast == nil {
			return nil
		}
		evts = tiles.StairsWestEast.fillRect(
			rng,
			image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y-1),
			a.cfg.ZOffsetWaterfall,
			propertiesRoad,
		)
	case collisionWaterfallNS:
		if tiles.WaterfallNorthSouth == nil {
			return nil
		}
		evts = tiles.WaterfallNorthSouth.fillRect(rng, r, a.cfg.ZOffsetWaterfall, propertiesWFall)
	case collisionWaterfallSN:
		if tiles.WaterfallSouthNorth == nil {
			return nil
		}
		evts = tiles.WaterfallSouthNorth.fillRect(
			rng,
			image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y+1),
			a.cfg.ZOffsetWaterfall,
			propertiesWFall,
		)
	case collisionWaterfallEW:
		if tiles.WaterfallEastWest == nil {
			return nil
		}
		evts = tiles.WaterfallEastWest.fillRect(
			rng,
			image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y),
			a.cfg.ZOffsetWaterfall,
			propertiesWFall,
		)
	case collisionWaterfallWE:
		if tiles.WaterfallWestEast == nil {
			return nil
		}
		evts = tiles.WaterfallWestEast.fillRect(
			rng,
			image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y),
			a.cfg.ZOffsetWaterfall,
			propertiesWFall,
		)
	}

	return evts
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

func (a *Autotiler) placeLand(o Outline, rng *rand.Rand, me *area, tagsonly bool) ([]*Event, string, error) {
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
	if beach > 0 && me.Data.Height() < a.cfg.CliffLevel {
		nearWtr = len(withinRadius(o, me.X, me.Y, beach, func(in *area) bool { return in.Data.IsWater() })) > 0
		nearWtrPlus = len(withinRadius(o, me.X, me.Y, beach+1, func(in *area) bool { return in.Data.IsWater() })) > 0
	}
	tsn := a.cfg.TransitionWidth

	if me.Data.Temperature() <= a.cfg.SnowLevel-tsn {
		src = firstFull(rng, tiles.Snow, tiles.Dirt, tiles.Rock)
		tag = Snow
	} else if nearWtr {
		src = firstFull(rng, tiles.Sand, tiles.Rock)
		tag = Sand
	} else if me.Data.Height() >= a.cfg.MountainLevel+tsn {
		src = firstFull(rng, tiles.Rock, tiles.Dirt)
		tag = Rock
	} else if me.Data.Temperature() <= a.cfg.VegetationMinTemp-tsn {
		src = firstFull(rng, tiles.Dirt, tiles.Rock)
		tag = Dirt
	} else if me.Data.Temperature() >= a.cfg.VegetationMaxTemp+tsn { // desert
		src = firstFull(rng, tiles.Sand, tiles.Rock, tiles.Dirt)
		tag = Sand
	}
	if tagsonly {
		return nil, tag, nil
	}

	if me.Data.Temperature() <= a.cfg.SnowLevel {
		srcT = firstTransition(rng, tiles.Snow, tiles.Dirt, tiles.Rock)
	} else if nearWtrPlus && !nearWtr {
		srcT = firstTransition(rng, tiles.Sand, tiles.Dirt)
	} else if me.Data.Temperature() >= a.cfg.VegetationMaxTemp {
		srcT = firstTransition(rng, tiles.Sand, tiles.Rock, tiles.Dirt)
	} else if me.Data.Temperature() <= a.cfg.VegetationMinTemp {
		srcT = firstTransition(rng, tiles.Dirt, tiles.Rock)
	} else if me.Data.Height() >= a.cfg.MountainLevel {
		srcT = firstTransition(rng, tiles.Rock, tiles.Dirt)
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
	src := tiles.Water.choosePiece(rng, crd, func(a *area) bool { return a.Data.IsWater() })

	evts := []*Event{
		newEvent(me.X, me.Y, a.cfg.ZOffsetWater, src, propertiesWater),
	}

	if tiles.Bridge == nil || !me.Data.IsRoad() {
		return evts, Water, nil
	}

	bsrc := tiles.Bridge.choosePiece(rng, crd, func(a *area) bool { return a.Data.IsRoad() && a.Data.IsWater() })
	evts = append(evts, newEvent(me.X, me.Y, a.cfg.ZOffsetRoad, bsrc, propertiesRoad))

	return evts, Road, nil
}

func (a *Autotiler) placeRoad(o Outline, rng *rand.Rand, me *area, tagonly bool) ([]*Event, string, error) {
	if !me.Data.IsRoad() {
		return nil, "", nil
	}
	tiles := me.Data.Tiles()
	ts := tiles.Road
	if me.Data.IsWater() {
		ts = tiles.Bridge
	}
	if ts == nil {
		return nil, "", nil
	}
	if tagonly {
		return nil, Road, nil
	}

	crd := cardinals(o, me.X, me.Y)
	src := ts.choosePiece(rng, crd, func(a *area) bool { return a.Data.IsRoad() })

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

	crd := cardinals(o, me.X, me.Y)

	src := tiles.Lava.choosePiece(rng, crd, func(a *area) bool { return a.Data.IsMolten() && !a.Data.IsWater() })
	return []*Event{
		newEvent(me.X, me.Y, a.cfg.ZOffsetWater, src, propertiesLava),
	}, Lava, nil

}

func (a *Autotiler) placeCliffs(o Outline, rng *rand.Rand, me *area, tagonly bool) ([]*Event, string, error) {
	h := me.Data.Height()
	if h < a.cfg.CliffLevel {
		return nil, "", nil
	}
	tiles := me.Data.Tiles()
	if tiles.Cliff == nil {
		return nil, "", nil
	}

	crd := cardinals(o, me.X, me.Y)
	src := tiles.Cliff.choosePiece(rng, crd, func(a *area) bool { return a.Data.Height() >= h })
	if src == "" {
		return nil, "", nil
	}
	evts := []*Event{
		newEvent(me.X, me.Y, a.cfg.ZOffsetCliff, src, propertiesCliff),
	}
	if !(me.Data.IsRoad() || me.Data.IsWater()) {
		return evts, CliffFace, nil
	}

	n := crd.North.Data.Height()
	s := crd.South.Data.Height()
	e := crd.East.Data.Height()
	w := crd.West.Data.Height()

	if me.Data.IsWater() {
		var wfalltype collisionType

		wrNS := crd.North.Data.IsWater() && crd.South.Data.IsWater()
		wrWE := crd.East.Data.IsWater() && crd.West.Data.IsWater()

		switch {
		case wrNS && n != s && n >= h && s <= h:
			wfalltype = collisionWaterfallNS
		case wrNS && n != s && n <= h && s >= h:
			wfalltype = collisionWaterfallSN
		case wrWE && w != e && e >= h && w <= h:
			wfalltype = collisionWaterfallEW
		case wrWE && w != e && e <= h && w >= h:
			wfalltype = collisionWaterfallWE
		}

		if wfalltype != "" {
			evts = append(evts, newColEvent(me.X, me.Y, a.cfg.ZOffsetCliff, wfalltype))
		}
	} else if me.Data.IsRoad() {
		var stype collisionType

		rdNS := crd.North.Data.IsRoad() && crd.South.Data.IsRoad()
		rdEW := crd.West.Data.IsRoad() && crd.East.Data.IsRoad()

		switch {
		case rdNS && n != s && n >= h && s <= h:
			stype = collisionStairsNS
		case rdNS && n != s && n <= h && s >= h:
			stype = collisionStairsSN
		case rdEW && e != w && e >= h && w <= h:
			stype = collisionStairsEW
		case rdEW && e != w && e <= h && w >= h:
			stype = collisionStairsWE
		}

		if stype != "" {
			evts = append(evts, newColEvent(me.X, me.Y, a.cfg.ZOffsetCliff, stype))
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
