package autotile

import (
	"github.com/voidshard/tile"

	"fmt"
	"sync"
)

const (
	// Water represents *any* water tile (regardless of swamp/sea/river etc)
	Water string = "water"

	// Specific kinds of water
	Swamp string = "swamp"
	River string = "river"
	Sea   string = "sea"

	// Ground is a catch all for *any* "solid ground" ie. not water / lava
	Ground string = "ground"

	// Specific kinds of ground
	Grass string = "grass"
	Sand  string = "sand"
	Dirt  string = "dirt"
	Snow  string = "snow"
	Rock  string = "rock"

	// Other misc options
	Road      string = "road"
	Lava      string = "lava"
	CliffFace string = "cliff"
)

// ObjectBin holds objects of varying types & handles choosing randomly via weighted
// chances & tags.
type ObjectBin struct {
	// reference to a user given loader
	load Loader

	// map object name -> tile pap
	objects map[string]*tile.Map

	// map group name -> chance object is placed
	chances map[string]float64

	// map object name -> object group
	objgroup map[string]string

	// map object group -> object names
	groups map[string][]string

	// tagsAll maps group -> tags (all)
	tagsAll map[string][]string

	// tagsAny maps group -> tags (any)
	tagsAny map[string][]string

	// how likely it is that we place nothing
	nilChance float64
}

// NewObjectBin creates a new ObjectBin that loads map via the given loader
func NewObjectBin(ldr Loader) *ObjectBin {
	return &ObjectBin{
		load:     ldr,
		objects:  map[string]*tile.Map{},
		chances:  map[string]float64{},
		objgroup: map[string]string{},
		groups:   map[string][]string{},
		tagsAll:  map[string][]string{},
		tagsAny:  map[string][]string{},
	}
}

// choose picks one of the given named objects considering their weights / tags for
// the given location (x, y, z).
func (o *ObjectBin) choose(mo *MapOutline, x, y, z int) *tile.Map {
	// firstly, check for a nil roll, since that vastly cuts down on our work
	n := mo.rng.Float64()
	if n <= o.nilChance {
		return nil // we rolled a `place nothing here`
	}

	sofar := o.nilChance
	for group, chance := range o.chances {
		// run through the groups, checking if we rolled that group
		if chance <= 0 {
			continue
		}
		sofar += chance
		if n > sofar {
			continue
		}
		// else: implies n <= sofar

		all, _ := o.tagsAll[group]
		any, _ := o.tagsAny[group]

		names, ok := o.groups[group]
		if !ok {
			continue
		}
		if len(names) == 0 {
			if matchTags(mo, all, any, x, y) {
				return nil
			} else {
				continue
			}
		}

		pickable := []*tile.Map{}
		for _, name := range names {
			// we want an obj from this group, but we can only pick objects
			// that fit & have their base match our tags.
			obj, ok := o.objects[name]
			if !ok {
				continue
			}

			// object doesn't fit without overwriting existing tiles -> never place
			if !mo.Tilemap.Fits(x, y, z, obj) {
				continue
			}

			// check that the base (bottom layer) of object sits on tiles
			// with matching tags.
			objheight := baseHeight(obj)
			suitable := true
			for ty := y + obj.Height - objheight; ty < y+obj.Height; ty++ {
				for tx := x; tx < x+obj.Width; tx++ {
					suitable = matchTags(mo, all, any, tx, ty)
				}
				if !suitable {
					break
				}
			}
			if suitable {
				pickable = append(pickable, obj)
			}
		}

		switch len(pickable) {
		case 0:
			continue
		case 1:
			return pickable[0]
		default:
			return pickable[mo.rng.Intn(len(pickable))]
		}
	}

	return nil
}

// matchTags returns of the given tile (x, y) has all tags in 'all'
// and at least one of the tags in 'any'
// Passing nils / no tags causes us to consider that check 'true'
// That is matchTags(outline, nil, nil, x, y) => true
func matchTags(mo *MapOutline, all, any []string, x, y int) bool {
	if all != nil {
		for _, t := range all {
			if !mo.HasTag(x, y, t) {
				return false
			}
		}
	}
	if any == nil || len(any) == 0 {
		return true
	}
	return mo.HasAnyTags(x, y, any)
}

// obj is an internal struct used during loading
type obj struct {
	Name string
	Map  *tile.Map
}

// Load a group of objects ('tob' .tmx files) & set their internal chance & tags.
//
// `group` is a name for the group (giving the same name twice will overwrite previous group).
// `chance` a base chance (0-1) for placing an object from this list.
// `objects` a list of object keys, these are passed to the `Loader` interface for retrieval.
// `all` is a list of tags base tiles must have in order to place one of the group
// `any` is a list of tags base tiles should have at least one of in order to place one of the group
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
//  - objects are loaded in parallel so the loader is required to be thread-safe.
//  - objects without tags are considered placable on any tiles
//  - in the same way a nil `all` tags or `any` tags implies that we're happy with anything
//  - objects will never be placed if they would overwrite existing tiles (regardless of tags)
//  - chances are not normalised, you'll probably want to manage this yourself
//  - if we're specifically given a group with no objects we will not place objects on
//    tiles with matching tags (if rolled)
func (o *ObjectBin) Load(group string, chance float64, objects, all, any []string) error {
	if group == "" {
		o.nilChance = chance
		return nil
	}
	if len(objects) == 0 {
		o.groups[group] = []string{}
		o.chances[group] = chance
		o.tagsAll[group] = all
		o.tagsAny[group] = any
		return nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(objects))
	objs := make(chan *obj)
	errs := make(chan error)

	for _, n := range objects {
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

	names := []string{}
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
		names = append(names, obj.Name)
		o.objects[obj.Name] = obj.Map
		o.objgroup[obj.Name] = group
	}
	o.groups[group] = names
	o.chances[group] = chance
	o.tagsAll[group] = all
	o.tagsAny[group] = any

	if !failed {
		final = nil
	}

	return final
}

// baseHeight returns the height in tiles of the lowest z-layer
// (from the bottom) in the given tile map (tob).
//
// so assuming a 3x5 map whose bottom most layer yields
// . . .
// . . .
// . x .
// x x x
// x x x
// where . is nil, x is a non-nil tile, we'd be hoping for `3`
func baseHeight(m *tile.Map) int {
	layers := m.ZLevels()
	if layers == nil || len(layers) == 0 {
		return 0
	}
	first := layers[0]

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.At(x, y, first) == nil {
				continue
			}

			return m.Height - y
		}
	}

	return 0 // there are no tiles on this layer so :shrug:
}
