package autotile

import (
	"math"
	"math/rand"
	"sync"

	"github.com/voidshard/tile"
)

// Autotiler is a struct that understands how to place various sets of tiles in in order
// to create a tiled map.
// - Only "orthogonal" square tiled maps are supported.
type Autotiler struct {
	cfg *Config
}

// NewAutotiler creates & returns an autotiler object.
func NewAutotiler(cfg *Config) (*Autotiler, error) {
	at := &Autotiler{cfg: cfg}
	return at, cfg.Validate()
}

// AddObjects adds larger objects (see the `tile` library on `tob` "tile objects")
// from an object bin to a map. This process is random-esqe but influenced by
// probabilities & tags (see objectbin.Load())
//
// An object bin holds sets of unique objects of various sizes belonging to groups.
// Each group of objects has a chance of being selected (that is, the chance that
// we place an object from that group) and a set of tags.
//
// Tags control where (in general terms) objects are allowed to be placed,
// in that their base (lowest z-layer) of tiles must sit on tiles matching the
// given group tags (all of `all` and at least one of `any`).
// We set tags for things like water, sea, river, snow, grass, rock, ground etc,
// but the user can also set tags via having their Outline interface return Area
// structs with `Tags` set (these tags are appended to internally set tags during
// the CreateMaps() function).
//
// In addition the object bin has a "nil" chance -- the probability that we do not
// place an object at a given location.
//
// Note
// - we never place objects if *any* of their tiles would overwrite an
// existing tile on the given map.
// - object bins are not threadsafe, AddObjects() calls should not share bins.
func (a *Autotiler) AddObjects(mo *MapOutline, bin *ObjectBin) {
	bin.normalise()

	for y := 0; y < mo.Tilemap.Height; y++ {
		for x := 0; x < mo.Tilemap.Width; x++ {
			obj := bin.choose(mo, x, y, a.cfg.ZOffsetObject)
			if obj == nil {
				continue
			}
			mo.Tilemap.Add(x, y, a.cfg.ZOffsetObject, obj)
		}
	}
}

// CreateMaps performs an initial tiling to create the base maps.
// Note that we're placing tiles one by one rather than tile-objects (that comes later).
// We place
// - water (rivers, swamp, sea)
// - grass / dirt / rock / sand / snow (underfoot)
// - cliffs
// - waterfalls
// - lava
// - roads
// We return completed base maps as they're done in order for the caller to perform
// any additional logic & ultimately write them out to disk.
// This is desirable as holding a large number of maps in memory consumes a massive
// amount of memory. When all maps are completed the returned channel is closed.
//
// This func is intended to enlarge a more distant world map (given by `Outline`)
// into a series of 'zoomed in' area maps.
func (a *Autotiler) CreateMaps(o Outline) (<-chan *MapOutline, error) {
	work := make(chan [2]int)
	wg := &sync.WaitGroup{}

	out := make(chan *MapOutline)

	for i := 0; i < a.cfg.Routines; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for coords := range work {
				out <- a.createMap(o, coords[0], coords[1])
			}
		}()
	}

	go func() {
		// spawn a routine to close the chan when all routines
		// have exited to let the caller know we're done
		wg.Wait()
		close(out)
	}()

	go func() {
		bnds := o.Bounds()
		h := bnds.Max.Y - bnds.Min.Y
		w := bnds.Max.X - bnds.Min.X

		for dy := 0; dy < h; dy += a.cfg.StepSize {
			for dx := 0; dx < w; dx += a.cfg.StepSize {
				work <- [2]int{dx / a.cfg.StepSize, dy / a.cfg.StepSize}
			}
		}

		close(work)
	}()

	return out, nil
}

// createMap makes a single base map offset at (mx, my) in world co-ords.
func (a *Autotiler) createMap(o Outline, mx, my int) *MapOutline {
	// build metadata for map
	meta := a.newMapOutline(o, mx, my)

	// place base tiles
	for ty := 0; ty < meta.MapHeight; ty++ {
		for tx := 0; tx < meta.MapWidth; tx++ {
			me := a.AtMapCoord(o, mx, my, tx, ty)

			if meta.numland > 0 {
				a.placeLand(o, meta, me)
			}

			if meta.numwater > 0 {
				a.placeWater(o, meta, me)
			}

			if meta.numlava > 0 {
				a.placeLava(o, meta, me)
			}

			if meta.numroad > 0 {
				a.placeRoad(o, meta, me)
			}

			if meta.numhighlands > 0 {
				a.placeCliffs(o, meta, me)
			}

			if me.Tags != nil {
				// add user custom tags
				tags := append(meta.Tags(me.X, me.Y), me.Tags...)
				meta.SetTags(me.X, me.Y, tags...)
			}
		}
	}

	return meta
}

// placeRoad handles laying down road tiles
func (a *Autotiler) placeRoad(o Outline, meta *MapOutline, me *Area) {
	if !(me.Road || meta.IsRoad(me.X, me.Y)) {
		return
	}
	if me.Lava || me.isWater() || meta.IsWater(me.X, me.Y) {
		return // TODO: bridge?
	}

	crd := a.cardinals(o, meta.MapX, meta.MapY, me.X, me.Y)
	in := []*Area{}
	out := []*Area{}
	for _, t := range crd.all() {
		if t.Road || meta.IsRoad(t.X, t.Y) {
			in = append(in, t)
		} else {
			out = append(out, t)
		}
	}

	// find required tiles
	tiles := me.Land
	if tiles.Road == nil {
		return
	}

	// choose the correct piece
	src := tiles.Road.choosePiece(meta.rng, in, out)
	meta.Tilemap.Set(me.X, me.Y, a.cfg.ZOffsetRoad, src)
	meta.Tilemap.SetProperties(src, propertiesRoad)
	meta.SetTags(me.X, me.Y, Road)
}

// placeCliffs handles laying down cliffs & waterfalls
func (a *Autotiler) placeCliffs(o Outline, meta *MapOutline, me *Area) {
	if me.Height < a.cfg.Params.CliffLevel {
		return
	}

	tiles := me.Land
	if tiles.Cliff == nil {
		return
	}

	crd := a.cardinals(o, meta.MapX, meta.MapY, me.X, me.Y)
	lowland := crd.Lower(me.Height)

	cliffs := tiles.Cliff.placements(meta.rng, me, lowland)
	for _, cp := range cliffs {
		meta.Tilemap.Set(cp.X, cp.Y, a.cfg.ZOffsetCliff, cp.Src)
		meta.Tilemap.SetProperties(cp.Src, propertiesCliff)
		meta.SetTags(cp.X, cp.Y, CliffFace)
	}

	if tiles.Waterfall == nil {
		// we can't place a waterfall because we lack the tiles
		return
	}
	if !(me.isWater() || meta.IsWater(me.X, me.Y)) {
		// if we're not in water, then clearly we don't need a waterfall
		return
	}
	if len(lowland) <= 1 {
		// not a cliff, thus no need for a waterfall
		return
	}
	low := headings(lowland)
	if includes(low, cornerNE...) || includes(low, cornerNW...) || includes(low, cornerSE...) || includes(low, cornerSW...) {
		// we don't place waterfalls on cliff corners
		return
	}

	watertiles := []*Area{}
	for _, t := range crd.all() {
		if t.isWater() || meta.IsWater(t.X, t.Y) {
			watertiles = append(watertiles, t)
		}
	}
	if len(watertiles) < 5 { // TODO: I think this should be < 6
		return
	}

	if includes(low, North) { // waterfall flowing north
		if tiles.Waterfall.SN == nil {
			return
		}
		if len(watertiles) < 8 {
			return
		}
		// we can only see the top of the waterfall falling out of view
		src := one(meta.rng, tiles.Waterfall.SN.MidTop)
		meta.Tilemap.Set(me.X, me.Y, a.cfg.ZOffsetWaterfall, src)
		meta.Tilemap.SetProperties(src, propertiesWFall)
	} else if includes(low, South) { // waterfall flowing south
		if tiles.Waterfall.NS == nil {
			return
		}
		// check the tile above (Y-1) of this tile
		above := a.AtMapCoord(o, meta.MapX, meta.MapY, me.X, me.Y-1)
		abovecrd := a.cardinals(o, meta.MapX, meta.MapY, me.X, me.Y-1)
		abovelowland := abovecrd.Lower(above.Height)
		iswfalltop := len(abovelowland) < 2
		wet := headings(watertiles)

		srcs := []string{}
		if len(watertiles) >= 8 { // middle
			if iswfalltop {
				srcs = append(srcs, one(meta.rng, tiles.Waterfall.NS.MidTop))
				meta.Tilemap.Set(me.X, me.Y-1, a.cfg.ZOffsetWaterfall, srcs[len(srcs)-1])
			}
			srcs = append(srcs, one(meta.rng, tiles.Waterfall.NS.MidCentre))
			meta.Tilemap.Set(me.X, me.Y, a.cfg.ZOffsetWaterfall, srcs[len(srcs)-1])
			srcs = append(srcs, one(meta.rng, tiles.Waterfall.NS.MidBottom))
			meta.Tilemap.Set(me.X, me.Y+1, a.cfg.ZOffsetWaterfall, srcs[len(srcs)-1])
		} else if includes(wet, cornerNE...) && includes(wet, cornerSE...) { // left
			if iswfalltop {
				srcs = append(srcs, one(meta.rng, tiles.Waterfall.NS.LeftTop))
				meta.Tilemap.Set(me.X, me.Y-1, a.cfg.ZOffsetWaterfall, srcs[len(srcs)-1])
			}
			srcs = append(srcs, one(meta.rng, tiles.Waterfall.NS.LeftCentre))
			meta.Tilemap.Set(me.X, me.Y, a.cfg.ZOffsetWaterfall, srcs[len(srcs)-1])
			srcs = append(srcs, one(meta.rng, tiles.Waterfall.NS.LeftBottom))
			meta.Tilemap.Set(me.X, me.Y+1, a.cfg.ZOffsetWaterfall, srcs[len(srcs)-1])
		} else if includes(wet, cornerNW...) && includes(wet, cornerSW...) { // right
			if iswfalltop {
				srcs = append(srcs, one(meta.rng, tiles.Waterfall.NS.RightTop))
				meta.Tilemap.Set(me.X, me.Y-1, a.cfg.ZOffsetWaterfall, srcs[len(srcs)-1])
			}
			srcs = append(srcs, one(meta.rng, tiles.Waterfall.NS.RightCentre))
			meta.Tilemap.Set(me.X, me.Y, a.cfg.ZOffsetWaterfall, srcs[len(srcs)-1])
			srcs = append(srcs, one(meta.rng, tiles.Waterfall.NS.RightBottom))
			meta.Tilemap.Set(me.X, me.Y+1, a.cfg.ZOffsetWaterfall, srcs[len(srcs)-1])
		}
		for _, s := range srcs {
			meta.Tilemap.SetProperties(s, propertiesWFall)
		}
	}
}

// placeLava lays down lava/molten rock tiles
func (a *Autotiler) placeLava(o Outline, meta *MapOutline, me *Area) {
	if !me.Lava {
		return
	}
	if me.isWater() || meta.IsWater(me.X, me.Y) {
		return // we don't lay lava in water
	}

	landtiles := me.Land
	if landtiles == nil || landtiles.Lava == nil {
		return
	}

	molten := []*Area{}
	notmolten := []*Area{}

	crd := a.cardinals(o, meta.MapX, meta.MapY, me.X, me.Y)
	for _, t := range crd.all() {
		if t.Lava || meta.IsLava(t.X, t.Y) {
			molten = append(molten, t)
		} else {
			notmolten = append(notmolten, t)
		}
	}

	// choose the correct piece
	src := landtiles.Lava.choosePiece(meta.rng, molten, notmolten)
	meta.Tilemap.Set(me.X, me.Y, a.cfg.ZOffsetWater, src)
	meta.Tilemap.SetProperties(src, propertiesWater)
	meta.SetTags(me.X, me.Y, Lava)
}

// placeWater handles setting water tiles (rivers, swamp, sea)
func (a *Autotiler) placeWater(o Outline, meta *MapOutline, me *Area) {
	if !(me.isWater() || meta.IsWater(me.X, me.Y)) {
		return
	}

	// look up nearby tiles & sort into water / not water
	crd := a.cardinals(o, meta.MapX, meta.MapY, me.X, me.Y)
	in := []*Area{}
	out := []*Area{}
	for _, t := range crd.all() {
		if t.isWater() || meta.IsWater(t.X, t.Y) {
			in = append(in, t)
		} else {
			out = append(out, t)
		}
	}

	// find required tiles
	tiles := me.Land
	if tiles == nil || tiles.Water == nil {
		return
	}

	// choose the correct piece
	src := tiles.Water.choosePiece(meta.rng, in, out)
	meta.Tilemap.Set(me.X, me.Y, a.cfg.ZOffsetWater, src)
	meta.Tilemap.SetProperties(src, propertiesWater)

	// set tags
	if me.Sea {
		meta.SetTags(me.X, me.Y, Water, Sea)
	} else if me.River {
		meta.SetTags(me.X, me.Y, Water, River)
	} else if me.Swamp {
		meta.SetTags(me.X, me.Y, Water, Swamp)
	} else {
		meta.SetTags(me.X, me.Y, Water)
	}
}

// placeLand decides what base land tile(s) should be used in our map
func (a *Autotiler) placeLand(o Outline, meta *MapOutline, me *Area) {
	landtiles := me.Land
	if landtiles == nil {
		return
	}
	landtiles.setTags()

	src, tag := firstFull(meta.rng, landtiles.Grass, landtiles.Dirt, landtiles.Rock)
	srcT := ""

	beach := a.cfg.Params.BeachWidth
	nearSea := false
	nearSeaPlus := false
	if beach > 0 {
		nearSea = len(a.withinRadius(o, meta.MapX, meta.MapY, me.X, me.Y, beach, func(in *Area) bool { return in.Sea })) > 0
		nearSeaPlus = len(a.withinRadius(o, meta.MapX, meta.MapY, me.X, me.Y, beach+1, func(in *Area) bool { return in.Sea })) > 0
	}

	if me.Temperature < a.cfg.Params.SnowLevel {
		src, tag = firstFull(meta.rng, landtiles.Snow, landtiles.Dirt, landtiles.Rock)
	} else if me.Temperature == a.cfg.Params.SnowLevel {
		srcT, _ = firstTransition(meta.rng, landtiles.Snow, landtiles.Dirt, landtiles.Rock)
	} else if me.Temperature < a.cfg.Params.VegetationMinTemp {
		src, tag = firstFull(meta.rng, landtiles.Dirt, landtiles.Rock)
	} else if me.Temperature == a.cfg.Params.VegetationMinTemp {
		srcT, _ = firstTransition(meta.rng, landtiles.Dirt, landtiles.Rock)
	} else if nearSeaPlus && !nearSea {
		srcT, _ = firstTransition(meta.rng, landtiles.Sand, landtiles.Rock)
	} else if nearSea {
		src, tag = firstFull(meta.rng, landtiles.Sand, landtiles.Rock)
	} else if me.Height > a.cfg.Params.MountainLevel {
		src, tag = firstFull(meta.rng, landtiles.Rock, landtiles.Dirt)
	} else if me.Height == a.cfg.Params.MountainLevel {
		srcT, _ = firstTransition(meta.rng, landtiles.Rock, landtiles.Dirt)
	} else if me.Temperature > a.cfg.Params.VegetationMaxTemp+1 { // desert
		src, tag = firstFull(meta.rng, landtiles.Sand, landtiles.Rock, landtiles.Dirt)
	} else if me.Temperature == a.cfg.Params.VegetationMaxTemp+1 {
		srcT, _ = firstTransition(meta.rng, landtiles.Sand, landtiles.Rock, landtiles.Dirt)
	}

	if srcT != "" {
		meta.Tilemap.Set(me.X, me.Y, a.cfg.ZOffsetLand+1, srcT)
		meta.Tilemap.SetProperties(srcT, propertiesLand)
	}
	meta.SetTags(me.X, me.Y, Ground, tag)
	meta.Tilemap.Set(me.X, me.Y, a.cfg.ZOffsetLand, src)
	meta.Tilemap.SetProperties(src, propertiesLand)
}

// newMapOutline prepares an outline for a given map. This is our wrapper struct around
// metadata, the tilemap (*tile.Map) and some other tidbits useful in book keeping.
func (a *Autotiler) newMapOutline(o Outline, mx, my int) *MapOutline {
	tmap := tile.New(&tile.Config{
		TileWidth:  uint(a.cfg.TileSize),
		TileHeight: uint(a.cfg.TileSize),
		MapWidth:   uint(a.cfg.MapSize),
		MapHeight:  uint(a.cfg.MapSize),
	})

	props := tmap.MapProperties()
	props.SetInt("worldoffset-x0", mx)
	props.SetInt("worldoffset-y0", my)
	tmap.SetMapProperties(props)

	seed := a.cfg.Seed + int64(mx) - (int64(my) * int64(my))
	meta := &MapOutline{
		parent:    a,
		Tilemap:   tmap,
		flood:     make([]bool, tmap.Height*tmap.Width),
		road:      make([]bool, tmap.Height*tmap.Width),
		lava:      make([]bool, tmap.Height*tmap.Width),
		tags:      make([]map[string]bool, tmap.Height*tmap.Width),
		MapX:      mx,
		MapY:      my,
		MapWidth:  tmap.Width,
		MapHeight: tmap.Height,
		seed:      seed,
		rng:       rand.New(rand.NewSource(seed)),
	}

	// first, collect metadata about our area, we do this to avoid any work we can
	// in future .. if possible.
	for ty := 0; ty < tmap.Height; ty++ {
		for tx := 0; tx < tmap.Width; tx++ {
			me := a.AtMapCoord(o, mx, my, tx, ty)

			if me.Height >= a.cfg.Params.CliffLevel {
				meta.numhighlands++
			}

			if me.isWater() {
				meta.numwater++
			} else {
				meta.numland++
			}

			if me.Road {
				meta.numroad++
			}

			if me.Lava {
				meta.numlava++
			}
		}
	}

	if meta.numroad == 0 && meta.numwater == 0 && meta.numlava == 0 {
		// we don't need to do any flood calcs -> joy!
		return meta
	}

	for ty := 0; ty < tmap.Height; ty++ {
		for tx := 0; tx < tmap.Width; tx++ {
			me := a.AtMapCoord(o, mx, my, tx, ty)
			crd := a.cardinals(o, mx, my, tx, ty)

			if meta.numlava > 0 {
				if a.shouldFlood(me, crd, func(in *Area) bool { return in.Lava }) {
					for _, t := range a.withinRadius(o, mx, my, tx, ty, 2, func(in *Area) bool {
						return !in.Lava
					}) {
						meta.setLava(t.X, t.Y)
					}
				}
			}

			if meta.numwater > 0 {
				if a.shouldFlood(me, crd, func(in *Area) bool { return in.isWater() || meta.IsWater(in.X, in.Y) }) {
					for _, t := range a.withinRadius(o, mx, my, tx, ty, 2, func(in *Area) bool {
						return !in.isWater()
					}) {
						meta.setWater(t.X, t.Y)
					}
				}
			}

			if meta.numroad > 0 {
				if a.shouldFlood(me, crd, func(in *Area) bool { return in.Road }) {
					for _, t := range a.withinRadius(o, mx, my, tx, ty, 2, func(in *Area) bool {
						return !in.Road
					}) {
						meta.setRoad(t.X, t.Y)
					}
				}
			}
		}
	}

	return meta
}

// shouldFlood returns true if
// - the given tile is part of some complex type set
// - two opposite corners of the tile part of some complex set
// In such cases we probably want to designate this tile & nearby tiles are part of the
// complex set. Essentially this prevents harsh diagonals from making strange breaks in our
// complex set.
func (a *Autotiler) shouldFlood(me *Area, crd *nearby, inset func(in *Area) bool) bool {
	if inset(me) {
		return false
	}

	in := []*Area{}
	out := []*Area{}
	for _, t := range crd.all() {
		if inset(t) {
			in = append(in, t)
		} else {
			out = append(out, t)
		}
	}

	if len(in) != 4 || len(out) != 4 {
		return false
	}

	outtiles := headings(out)

	if includes(outtiles, cornerNE...) && includes(outtiles, SouthWest) {
		return true
	} else if includes(outtiles, cornerSE...) && includes(outtiles, NorthWest) {
		return true
	} else if includes(outtiles, cornerSW...) && includes(outtiles, NorthEast) {
		return true
	} else if includes(outtiles, cornerNE...) && includes(outtiles, SouthWest) {
		return true
	}
	return false
}

// AtMapCoord coverts a coord from map space to world space & returns the resulting
// area information.
// Ie. given the map outline o at map (mapx, mapy) what is at tile (tilex, tiley)
func (a *Autotiler) AtMapCoord(o Outline, mapx, mapy, tilex, tiley int) *Area {
	// return the Area from the landscape for the given tile (x,y) on this map
	finalx := (mapx * a.cfg.StepSize) + int(math.Floor((float64(tilex) * float64(a.cfg.StepSize) / float64(a.cfg.MapSize))))
	finaly := (mapy * a.cfg.StepSize) + int(math.Floor((float64(tiley) * float64(a.cfg.StepSize) / float64(a.cfg.MapSize))))

	r := o.At(finalx, finaly)
	if r == nil {
		return &Area{X: tilex, Y: tiley, Sea: true}
	}

	result := &Area{
		X:           tilex,
		Y:           tiley,
		Land:        r.Land,
		Height:      r.Height,
		Temperature: r.Temperature,
		Sea:         r.Sea,
		River:       r.River,
		Swamp:       r.Swamp,
		Lava:        r.Lava,
		Road:        r.Road,
		Tags:        r.Tags,
	}

	if result.Lava {
		// because placing cliffs in lava is weird
		result.Height = a.cfg.Params.CliffLevel - 1
	}

	return result
}
