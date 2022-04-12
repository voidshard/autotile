package draw

import (
	"github.com/voidshard/autotile"
	"github.com/voidshard/tile"

	"image"
	"math/rand"
	"time"
)

// FillRect fills a rect from min -> max with tiles using the given tileset pattern.
// - We fill up to and including the max.
// - We set the correct (we hope) tiles from the bottom up (ie. from max Y).
// - If tiles aren't set for a given piece we place nothing for that space.
//   Ie. if we needed a 3 quarter tile to place and we aren't given any then we simply
//   don't place a tile there & move on.
func FillRect(dst tile.Tileable, t *autotile.Tileset, in image.Rectangle, z int, props *tile.Properties) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	toSet := map[string]bool{} // properties will be set on these

	pickTile := func(tx, ty int, in []string) error {
		src := one(rng, in)
		if src == "" {
			return nil
		}
		toSet[src] = true
		return dst.Set(tx, ty, z, src)
	}

	for y := in.Max.Y; y >= in.Min.Y; y-- {
		for x := in.Min.X; x <= in.Max.X; x++ {
			var err error
			switch y {
			case in.Max.Y:
				switch x {
				case in.Min.X:
					err = pickTile(x, y, t.QuarterNorthEast)
				case in.Max.X:
					err = pickTile(x, y, t.QuarterNorthWest)
				default:
					err = pickTile(x, y, t.NorthHalf)
				}
			case in.Min.Y:
				switch x {
				case in.Min.X:
					err = pickTile(x, y, t.QuarterSouthEast)
				case in.Max.X:
					err = pickTile(x, y, t.QuarterSouthWest)
				default:
					err = pickTile(x, y, t.SouthHalf)
				}
			default:
				switch x {
				case in.Min.X:
					err = pickTile(x, y, t.EastHalf)
				case in.Max.X:
					err = pickTile(x, y, t.WestHalf)
				default:
					err = pickTile(x, y, t.Full)
				}
			}
			if err != nil {
				return err
			}
		}
	}

	if props != nil {
		for name := range toSet {
			err := dst.SetProperties(name, props)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// one chooses one item at random
func one(rng *rand.Rand, items []string) string {
	if items == nil || len(items) == 0 {
		return ""
	}
	return items[rng.Intn(len(items))]
}
