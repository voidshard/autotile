package estate

import (
	"fmt"
	"image"
	"math/rand"
	"strings"
	"time"

	"github.com/voidshard/tile"

	"github.com/voidshard/autotile"
	"github.com/voidshard/autotile/internal/binpack"
	"github.com/voidshard/autotile/pkg/draw"
)

// Estate represents some collections of objects w/ sub sets of objects
// bin packed into a neat area.
type Estate struct {
	cfg *Config
	rng *rand.Rand
	ldr autotile.Loader

	root   placements
	width  int
	height int
	offset int
}

// build kicks off our bin packing fun
func (e *Estate) build() error {
	root, err := newPlacements(e.cfg.Set, e.ldr, e.rng)
	if err != nil {
		return err
	}
	e.root = root

	w, h := binpack.Pack(root) // this is the floor area
	if e.cfg.Set.Fence != nil {
		w += 2
		h += 2
	}
	e.width = w
	e.height = h

	e.offset = e.calculateHeightOffset() // but things can extend above that
	return nil
}

// printString helps debugging
func (e *Estate) printString() {
	for _, child := range e.root {
		printString(child, 0)
	}
}

// printString helps debugging
func printString(p *placement, depth int) {
	fmt.Printf("%s%s %v %v\n", strings.Repeat("\t", depth), p.ptype, p.packLoc, p.packSize)
	if p.children != nil {
		for _, child := range p.children {
			printString(child, depth+1)
		}
	}
}

// calculateHeightOffset figures out how much higher the map needs to be
// in order to fit all the tiles above the lowest z level in.
// That is, for bin packing we only care about the obj footprint size
// (ie. a tree with 3 z levels but whose trunk is 1x1 is binpacked as 1x1)
// so this func works backwards to figure out if we need to shunt everything
// down a bit in order to fit in stuff higher than ground level (ie. tree canopy).
func (e *Estate) calculateHeightOffset() int {
	small := 0
	for _, child := range e.root {
		min := child.minY(image.Pt(0, 0))
		if min < small {
			small = min
		}
	}
	if small < 0 {
		return -1 * small
	}
	return 0
}

// Map turns our estate into a tile.Map for rendering / further processing
func (e *Estate) Map(tilewidth, tileheight uint) (*tile.Map, error) {
	tmap := tile.New(&tile.Config{
		MapWidth:   uint(e.width),
		MapHeight:  uint(e.height + e.offset),
		TileWidth:  tilewidth,
		TileHeight: tileheight,
	})

	depth := 0
	bnds := image.Rect(0, e.offset, e.width-1, e.height+e.offset-1)
	if e.cfg.Set.Base != nil {
		err := draw.FillRect(tmap, e.cfg.Set.Base, bnds, depth, e.cfg.Set.BaseProps)
		if err != nil {
			return nil, err
		}
		depth++
	}
	if e.cfg.Set.Fence != nil {
		err := draw.FillRect(tmap, e.cfg.Set.Fence, bnds, depth, e.cfg.Set.FenceProps)
		if err != nil {
			return nil, err
		}
		if e.cfg.Set.GateLocation == CentreBottom {
			err := tmap.Set((bnds.Max.X-bnds.Min.X)/2, bnds.Max.Y, depth, e.cfg.Set.Gate)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, p := range e.root {
		err := p.set(tmap, e.rng, image.Pt(0, e.offset), depth)
		if err != nil {
			return nil, err
		}
	}

	return tmap, nil
}

// Build an Estate from a config
func Build(cfg *Config) (*Estate, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config must be specified")
	}
	if cfg.Set == nil {
		return nil, fmt.Errorf("estate config must contain a Set struct")
	}

	seed := cfg.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	ldr := cfg.Loader
	if ldr == nil {
		ldr = autotile.NewFileLoader("")
	}

	est := &Estate{
		cfg: cfg,
		rng: rand.New(rand.NewSource(seed)),
		ldr: ldr,
	}

	return est, est.build()
}
