package autotile

import (
	perlin "github.com/voidshard/autotile/internal/perlin"
	"github.com/voidshard/tile"

	"fmt"
	"image"
	"math/rand"
	"sort"
	"sync"
)

// Distribution indicates how we'll distribute objects, or phrased another
// way; how we're going to generate random numbers to determine what objects
// go where.
// RandomDistribution is the default.
type Distribution string

const (
	// RandomDistribution means we use a normal RNG from the random Go package.
	RandomDistribution Distribution = "random"

	// PerlinDistribution uses a perlin noise function to determine random numbers
	// between 0-1 for the map area. This has the effect of distributing things
	// randomly but them appearing a little more ordered - patches of trees,
	// paths through forests etc.
	PerlinDistribution Distribution = "perlin"
)

// Bin holds objects of varying types & handles choosing randomly via weighted
// chances & tags.
// Nb. Bin is not considered threadsafe.
type Bin struct {
	// reference to a user given loader
	load Loader

	// map object name -> tile map
	objects map[string]*tile.Map

	// all loaded groups & their config
	groups     map[string]*BinGroupConfig
	groupOrder []string // order which we'll iterate during Choose)

	// how likely we are to place objects via a given model
	distChance map[Distribution]float64

	// how likely it is that we place nothing
	nilChance float64

	// perlinmap if distribution is PerlinDistribution
	perlinmap *image.RGBA

	// random number generator
	rng  *rand.Rand
	seed int64

	tiler      *Autotiler
	mapoutline Outline
}

// NewBin creates a new Bin that loads map via the given loader
func NewBin(a *Autotiler, o Outline, seed int64, ldr Loader) *Bin {
	return &Bin{
		tiler:      a,
		mapoutline: o,
		seed:       seed,
		rng:        rand.New(rand.NewSource(seed)),
		load:       ldr,
		objects:    map[string]*tile.Map{},
		groups:     map[string]*BinGroupConfig{},
	}
}

// normalise ensures our object groups have normalied probabilities within their
// distribution type & calculates the chance(s) that we place something in given
// distribution types.
func (o *Bin) normalise() {
	distTtl := map[Distribution]float64{
		RandomDistribution: 0.0,
		PerlinDistribution: 0.0,
	}

	allTotal := o.nilChance

	// normalised probabilities within each distribution type
	for _, g := range o.groups {
		total, _ := distTtl[g.Distribution]
		distTtl[g.Distribution] = total + g.Chance

		allTotal += total
	}
	for _, g := range o.groups {
		total, _ := distTtl[g.Distribution]
		g.normChance = g.Chance / total
	}

	// the chance we'll pick a group with a given distribution type
	dc := map[Distribution]float64{}
	for dist, distProb := range distTtl {
		dc[dist] = distProb / allTotal
	}

	o.distChance = dc
}

// SetDistribution sets how we'll generate random numbers
func (o *Bin) setPerlinDistribution() {
	o.perlinmap = perlin.New(1000, 1000, 0.06, o.seed)
}

// perlinValue yields a number 0-1 based on a perlin map value at (x,y)
func (o *Bin) perlinValue(x, y int) float64 {
	c := o.perlinmap.At(x%1000, y%1000)
	r, _, _, _ := c.RGBA() // it's greyscale anyways
	return float64(r) / 255 / 255
}

// Choose picks one of the given named objects considering their weights / tags for
// the given location (x, y, z).
//
// Essentially we need three random numbers
// - first a random number to determine what placement distribution we'll go with
// - secondly a random number generated according to that distribution to select which
//   group from that distribution to choose
// - thirdly a final random number to choose which of the placeable object from that
//   group to pick
//
// So assuming we had two groups with "PerlinDistribution" and two with "RandomDistribution"
// we first randomly decide either Perlin or Random.
// We then move on to picking a specific group from those within either Perlin or Random.
// Finally we determine what objects from the chosen group are placeable
// - if one is placeable -> done
// - if none are placeable -> move onto the next group
// - if more than one is placeable -> choose at random
func (o *Bin) Choose(t tile.Tileable, x, y, z int) (string, *tile.Map, error) {
	// nb. careful to iterate lists here & not dists (whose order is undefined)

	// firstly, check for a nil roll, since that vastly cuts down on our work
	rn := o.rng.Float64()
	if rn <= o.nilChance {
		return "", nil, nil // we rolled a `place nothing here`
	}

	// decide which distribution model we're going with
	var distributionModel Distribution
	sofar := o.nilChance
	for _, name := range []Distribution{RandomDistribution, PerlinDistribution} {
		chance, _ := o.distChance[name]

		if chance <= 0 {
			continue
		}

		sofar += chance
		if rn > sofar {
			continue
		}
		// else: implies rn <= sofar

		distributionModel = name

		switch distributionModel {
		case RandomDistribution:
			rn = o.rng.Float64()
		case PerlinDistribution:
			if o.perlinmap == nil {
				o.setPerlinDistribution()
			}
			rn = o.perlinValue(x, y)
		}

		break
	}

	// run through groups matching our chosen model & pick one
	sofar = 0.0
	for _, groupName := range o.groupOrder {
		cfg, ok := o.groups[groupName]
		if !ok || cfg.Distribution != distributionModel {
			continue // we've already chosen which kind we want
		}

		if cfg.normChance <= 0 {
			continue
		}
		sofar += cfg.normChance
		if rn > sofar {
			continue
		}
		// else: implies rn <= sofar

		// finally we need to check what specific objects of this group fit
		pickablenames := []string{}
		pickable := []*tile.Map{}
		for _, name := range cfg.Objects {
			// we want an obj from this group, but we can only pick objects
			// that fit & have their base match our tags.
			obj, ok := o.objects[name]
			if !ok {
				continue
			}

			// object doesn't fit without overwriting existing tiles -> never place
			fits, err := t.Fits(x, y, z, obj)
			if err != nil {
				return "", nil, err
			}
			if !fits {
				continue
			}

			// check that the base (bottom layer) of object sits on tiles
			// with matching tags.
			objheight := baseHeight(obj)
			suitable := true
			for ty := y + obj.Height - objheight; ty < y+obj.Height; ty++ {
				for tx := x; tx < x+obj.Width; tx++ {
					tiletags, err := o.tiler.TagsAt(o.mapoutline, tx, ty)
					if err != nil {
						return "", nil, err
					}

					suitable = matchTags(tiletags, cfg.TagsAll, cfg.TagsAny)
					if !suitable {
						break
					}
				}
			}
			if suitable {
				pickable = append(pickable, obj)
				pickablenames = append(pickablenames, name)
			}
		}

		switch len(pickable) {
		case 0:
			continue // nothing fits :(
		case 1:
			return pickablenames[0], pickable[0], nil
		default:
			num := o.rng.Intn(len(pickable))
			return pickablenames[num], pickable[num], nil
		}
	}

	return "", nil, nil
}

// matchTags returns if the given tile tags has all tags in 'all'
// and at least one of the tags in 'any'
// Passing nils / no tags causes us to consider that check 'true'
func matchTags(tiletags, all, any []string) bool {
	tags := map[string]bool{}
	for _, t := range tiletags {
		tags[t] = true
	}

	if all != nil {
		for _, t := range all {
			_, found := tags[t]
			if !found {
				return false
			}
		}
	}

	if any == nil || len(any) == 0 {
		return true
	}

	for _, t := range any {
		_, found := tags[t]
		if found {
			return true
		}
	}

	return false
}

// obj is an internal struct used during loading
type obj struct {
	Name string
	Map  *tile.Map
}

//
// `chance` a base chance (0-1) for placing an object from this list.
// `objects` a list of object keys, these are passed to the `Loader` interface for retrieval.
// `all` is a list of tags base tiles must have in order to place one of the group
// `any` is a list of tags base tiles should have at least one of in order to place one of the group
// `distribution` here indicates how objects of this group should be laid out, ie, how random numbers are generated
//
// [all | any] When considering whether we can place an obj based on it's tags, all of it's base (lowesr
// z-layer) tiles must fall on map tile(s) with matching tags.
// That is, if an object uses 10 z-layers, we'd only check that the bottom most z-layer (probably `0`)
// sits on matching tagged map tiles.
// By 'matching' we mean;
//  - each base tile must have all of the tags found in
// The phrase another way; assuming we were placing a building whose *map* size was 50x50
// but whose lowest z-layer was 50x10 (ie. the building's ground floor occupies 50x10) we only care
// that the ground floor tiles sit on tiles matching our requested tags.
//
// [Nil] As a special case an empty group name (ie "") will set the 'nil' object chance.
// That is, the chance that we deliberately place *no* object at all.
//
// Nb:
//  - objects without tags are considered placable on any tiles
//  - in the same way a nil `all` tags or `any` tags implies that we're happy with anything
//  - objects will never be placed if they would overwrite existing tiles (regardless of tags)
//  - if we're specifically given a group with no objects we will not place objects on
//    tiles with matching tags (if rolled)
type BinGroupConfig struct {
	// Chance is the probability that we will try to place an object from this group
	Chance float64

	// normChance is Chance normalised against other groups
	normChance float64

	// Objects is the list of objects belonging to this group. It is expected that
	// even between groups each name yields a unique object (.tob)
	Objects []string

	// TagsAll indicates that all base tiles of this object must have each of these tags
	TagsAll []string

	// TagsAny indicates that all base tiles of this object must have at least one
	// of these tags
	TagsAny []string

	// Distribution indicates how randomness is determined for this group
	Distribution Distribution
}

func (l *BinGroupConfig) applyDefaults() {
	if string(l.Distribution) == "" {
		l.Distribution = RandomDistribution
	}
}

// Load a group of objects ('tob' .tmx files) & set their internal chance & tags.
// Nb:
//  - objects are loaded in parallel so the loader is required to be thread-safe.
func (o *Bin) Load(group string, cfg *BinGroupConfig) error {
	cfg.applyDefaults()

	if group == "" {
		o.nilChance = cfg.Chance
		return nil
	}
	if cfg.Objects == nil || len(cfg.Objects) == 0 {
		return nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(cfg.Objects))
	objs := make(chan *obj)
	errs := make(chan error)

	for _, n := range cfg.Objects {
		go func(name string) {
			defer wg.Done()
			objmap, err := o.load.Map(name)
			if err != nil {
				errs <- err
			} else {
				objs <- &obj{
					Name: name,
					Map:  objmap,
				}
			}
		}(n)
	}

	go func() {
		wg.Wait()
		close(objs)
		close(errs)
	}()

	failed := false
	final := fmt.Errorf("failed to load map(s)")

	go func() {
		for err := range errs {
			// roll up errors
			failed = true
			final = fmt.Errorf("%w %v", final, err)
		}
	}()

	for obj := range objs {
		// insert into our internal maps
		o.objects[obj.Name] = obj.Map
	}
	o.groups[group] = cfg
	o.groupOrder = append(o.groupOrder, group)

	if !failed {
		final = nil
	}

	sort.Strings(o.groupOrder)
	o.normalise()
	return final
}

// baseHeight returns the height in tiles of the lowest z-layer
// (from the bottom) in the given tile map (tob).
//
// so assuming a 3x5 map whose bottom most layer yields
// . . .
// . . .
// . x .
// . x x
// x x x
// where . is nil, x is a non-nil tile, we'd be hoping for `3`
// since on the second column an x reaches the 3rd row, counting
// upwards from the bottom.
func baseHeight(m *tile.Map) int {
	// determine lowest layer
	layers := m.ZLevels()
	if layers == nil || len(layers) == 0 {
		return 0
	}
	first := layers[0]

	// walk from the top left corner until we find the first
	// non nil tile
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.At(x, y, first) == nil {
				continue
			}

			// since we want height from the bottom
			return m.Height - y
		}
	}

	return 0 // there are no tiles on this layer so :shrug:
}
