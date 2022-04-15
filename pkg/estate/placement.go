package estate

import (
	"image"
	"math/rand"

	"github.com/voidshard/tile"

	"github.com/voidshard/autotile"
	"github.com/voidshard/autotile/internal/binpack"
	"github.com/voidshard/autotile/pkg/draw"
)

// placeType differentiates between the different origins for
// placement structs
type placeType string

var (
	obj   placeType = "object"
	set   placeType = "set"
	space placeType = "space"
)

// placements implements Packer for internal/binpack
type placements []*placement

func (p placements) Len() int              { return len(p) }
func (p placements) Size(n int) (int, int) { return p[n].size() }
func (p placements) Place(n, x, y int) {
	if p[n].parentSet.Fence != nil {
		x++
		y++
	}
	if p[n].ptype != space {
		x += p[n].parentSet.PadLeft
		y += p[n].parentSet.PadTop
	}
	pnt := image.Pt(x, y)
	p[n].packLoc = &pnt
}

// placement is a single thing we want to put on the finished map
type placement struct {
	ptype    placeType
	children placements

	parentSet *Set

	targetObj *tile.Map
	targetSet *Set

	packSize *image.Point
	packLoc  *image.Point
}

// minY calculating the lowest Y value this placement needs to exist
// in order to fit the whole thing in the map.
func (p *placement) minY(off image.Point) int {
	at := image.Pt(off.X+p.packLoc.X, off.Y+p.packLoc.Y)

	switch p.ptype {
	case obj:
		return at.Y - p.targetObj.Height + autotile.BaseHeight(p.targetObj)
	case set:
		var min *int

		for _, child := range p.children {
			cm := child.minY(at)
			if min == nil || cm < *min {
				min = &cm
			}
		}

		return *min
	}

	return 0
}

// set draws this placements onto the Tileable `tmap` with offset `off`
func (p *placement) set(tmap tile.Tileable, rng *rand.Rand, off image.Point, depth int) error {
	at := image.Pt(off.X+p.packLoc.X, off.Y+p.packLoc.Y)

	switch p.ptype {
	case obj:
		err := tmap.Add(at.X, at.Y-p.targetObj.Height+autotile.BaseHeight(p.targetObj), depth, p.targetObj)
		if err != nil {
			return err
		}
	case set:
		areaExclPadding := image.Rect(
			at.X,
			at.Y,
			at.X+p.packSize.X-p.parentSet.PadRight-p.parentSet.PadLeft-1,
			at.Y+p.packSize.Y-p.parentSet.PadBottom-p.parentSet.PadTop-1,
		)
		if p.targetSet.Base != nil {
			err := draw.FillRect(tmap, p.targetSet.Base, areaExclPadding, depth, p.targetSet.BaseProps)
			if err != nil {
				return err
			}
			depth++
		}
		if p.targetSet.Fence != nil {
			err := draw.FillRect(tmap, p.targetSet.Fence, areaExclPadding, depth, p.targetSet.FenceProps)
			if err != nil {
				return err
			}
			if p.targetSet.GateLocation == CentreBottom {
				tmap.Set(
					areaExclPadding.Max.X-(areaExclPadding.Max.X-areaExclPadding.Min.X)/2,
					areaExclPadding.Max.Y,
					depth,
					p.targetSet.Gate,
				)
			}
		}
		for _, child := range p.children {
			err := child.set(tmap, rng, at, depth)
			if err != nil {
				return err
			}
		}
	case space:
		// TODO
	}

	return nil
}

// size returns the width, height of this item in tiles
func (p *placement) size() (int, int) {
	if p.packSize != nil {
		return p.packSize.X, p.packSize.Y
	}

	w, h := 0, 0

	switch p.ptype {
	case obj:
		// ie. return the floor area of the lowest layer
		w, h = p.targetObj.Width, autotile.BaseHeight(p.targetObj)
	case set:
		w, h = binpack.Pack(p.children)
		if p.targetSet.Fence != nil {
			w += 2
			h += 2
		}
	case space:
		// we always set packSize on creation of a `space`
	}

	w += p.parentSet.PadLeft + p.parentSet.PadRight
	h += p.parentSet.PadTop + p.parentSet.PadBottom

	pnt := image.Pt(w, h)
	p.packSize = &pnt
	return w, h
}

// between gets random number between min-max
func between(rng *rand.Rand, min, max int) int {
	if min > max {
		return 0
	}
	if min == max {
		return 1
	}
	return rng.Intn(max-min) + min
}

// newPlacements returns new placements for the given set `s`
// This implies we walk the tree downward calling newPlacements in
// turn for each set with children, yielding a complete tree.
func newPlacements(s *Set, ldr autotile.Loader, rng *rand.Rand) (placements, error) {
	me := []*placement{}

	// add objects
	if s.Objects != nil {
		for _, in := range s.Objects {
			num := between(rng, in.MinCount, in.MaxCount)
			if num <= 0 {
				continue
			}
			tmap, err := ldr.Map(in.Tob)
			if err != nil {
				return nil, err
			}

			for i := 0; i < num; i++ {
				pl := &placement{
					parentSet: s,
					ptype:     obj,
					targetObj: tmap,
				}

				w, h := pl.size()
				me = append(me, pl)

				// and some brief consideration of empty space
				if s.EmptyPercentage > 0 {
					ps := image.Point{}

					r := rng.Float64()
					if r <= 0.5 {
						ps.X = int(float64(w) * s.EmptyPercentage)
						ps.Y = h
					} else {
						ps.X = w
						ps.Y = int(float64(h) * s.EmptyPercentage)
					}

					pl := &placement{
						parentSet: s,
						ptype:     space,
						packSize:  &ps,
					}

					me = append(me, pl)
				}
			}
		}
	}

	// add our explicitly sized empty space(s)
	if s.Empty != nil {
		for _, in := range s.Empty {
			if in.X <= 0 || in.Y <= 0 {
				continue
			}
			me = append(me, &placement{
				parentSet: s,
				ptype:     space,
				packSize:  &in,
			})
		}
	}

	// run through our child sets
	if s.Sets != nil {
		for _, cs := range s.Sets {
			children, err := newPlacements(cs, ldr, rng)
			if err != nil {
				return nil, err
			}

			if len(children) == 0 {
				continue // since we contain nothing ..
			}
			pl := &placement{
				parentSet: s,
				ptype:     set,
				targetSet: cs,
				children:  children,
			}

			me = append(me, pl)
		}
	}

	return placements(me), nil
}
