package main

import (
	"log"

	"github.com/voidshard/autotile"
	"github.com/voidshard/tile"
)

var tobs = []struct {
	Group        string
	Chance       float64
	Names        []string
	All          []string
	Any          []string
	Distribution autotile.Distribution
}{
	{
		"",
		0.55, // 55% of the time, place nothing
		nil,
		nil,
		nil,
		autotile.RandomDistribution,
	},
	{ // place either of these objects 25% of the time, but only on dirt or grass tiles
		"grass",
		0.25,
		[]string{"grass.short.05.tmx", "grass.short.06.tmx"},
		nil,
		[]string{autotile.Dirt, autotile.Grass},
		autotile.RandomDistribution,
		//autotile.PerlinDistribution,
	},
	{ // place 2% of the time but only on sand
		"grass-on-sand",
		0.02,
		[]string{"grass.short.05.tmx"},
		nil,
		[]string{autotile.Sand},
		autotile.RandomDistribution,
		//autotile.PerlinDistribution,
	},
	{ // place 3% of the time on dirt or grass
		"shrooms",
		0.03,
		[]string{"mushroom.01.tmx", "mushroom.02.tmx", "mushroom.03.tmx"},
		nil,
		[]string{autotile.Dirt, autotile.Grass},
		autotile.RandomDistribution,
	},
	{ // place pretty much anywhere, but only 1% of the time
		"rocks",
		0.01,
		[]string{"standingrock.03.tmx", "standingrock.04.tmx"},
		nil,
		[]string{autotile.Dirt, autotile.Grass, autotile.Sand, autotile.Rock},
		autotile.RandomDistribution,
	},
	{ // place 15% of the time on dirt or grass
		"vegetation",
		0.15,
		[]string{"tree.large.06.tmx", "tree.large.07.tmx", "tree.small.07.tmx", "tree.small.08.tmx", "tree.small.09.tmx"},
		nil,
		[]string{autotile.Dirt, autotile.Grass},
		autotile.RandomDistribution,
		//autotile.PerlinDistribution,
	},
}

func main() {
	/* Test func renders out example scene(s) so one can eyeball logic changes.
	*  Reads from test/tobs/ folder and sets images paths as if one was inside test/tiles/
	*  Maps are .tmx compatible and readable via the tiled map editor (you'll want to move
	*  the .tmx file inside of test/tiles/ & open it with tiled.
	*  Tileset is embedded automatically into the map for all required tiles.
	 */

	// first, we need a map (tile.Tileable) to actually write to
	tmap := tile.New(&tile.Config{
		TileWidth:  32,
		TileHeight: 32,
		MapWidth:   uint(beach.max * beach.scale),
		MapHeight:  uint(beach.max * beach.scale),
	})

	at, err := autotile.NewAutotiler(cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		// autotiler outputs events about what it is placing where as it goes,
		// if you call Events() you *must* consume them.
		// Events are not stored, if you're not listening when the event occurs
		// then you miss out
		for e := range at.Events() {
			if e.ObjectID != "" {
				log.Printf(" (%d,%d,%d) -> placed object %s\n", e.X, e.Y, e.Z, e.ObjectID)
			} else {
				log.Printf(" (%d,%d,%d) -> placed tile %s\n", e.X, e.Y, e.Z, e.Src)
			}
		}
	}()

	err = at.SetLand(beach, beach.Bounds(), tmap)
	if err != nil {
		panic(err)
	}

	ldr := autotile.NewFileLoader("test/tobs/")
	bin := autotile.NewBin(at, beach, 123456789, ldr)

	for _, grp := range tobs {
		err = bin.Load(
			grp.Group,
			&autotile.BinGroupConfig{
				Chance:       grp.Chance,
				Objects:      grp.Names,
				TagsAll:      grp.All,
				TagsAny:      grp.Any,
				Distribution: grp.Distribution,
			},
		)
		if err != nil {
			panic(err)
		}
	}

	err = at.SetObjects(beach, beach.Bounds(), tmap, bin)
	if err != nil {
		panic(err)
	}

	err = tmap.WriteFile("maptest.01.tmx")
	if err != nil {
		panic(err)
	}

}
