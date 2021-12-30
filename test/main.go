package main

import (
	"fmt"

	"github.com/voidshard/autotile"
)

var tobs = []struct {
	Group  string
	Chance float64
	Names  []string
	All    []string
	Any    []string
}{
	{
		"",
		0.52, // 52% of the time, place nothing
		nil,
		nil,
		nil,
	},
	{
		"grass",
		0.3,
		[]string{"grass.short.05.tmx", "grass.short.06.tmx"},
		nil,
		[]string{autotile.Dirt, autotile.Grass, autotile.Sand},
	},
	{
		"shrooms",
		0.05,
		[]string{"mushroom.01.tmx", "mushroom.02.tmx", "mushroom.03.tmx"},
		nil,
		[]string{autotile.Dirt, autotile.Grass},
	},
	{
		"rocks",
		0.03,
		[]string{"standingrock.03.tmx", "standingrock.04.tmx"},
		nil,
		[]string{autotile.Ground},
	},
	{
		"vegetation",
		0.1,
		[]string{"tree.large.06.tmx", "tree.large.07.tmx", "tree.small.07.tmx", "tree.small.08.tmx", "tree.small.09.tmx"},
		nil,
		[]string{autotile.Dirt, autotile.Grass},
	},
}

func main() {
	/* Test func renders out example scene(s) so one can eyeball logic changes.
	*  Reads from test/tobs/ folder and sets images paths as if one was inside test/tiles/
	*  Maps are .tmx compatible and readable via the tiled map editor (you'll want to move
	*  the .tmx file inside of test/tiles/ & open it with tiled.
	*  Tileset is embedded automatically into the map for all required tiles.
	 */
	at, err := autotile.NewAutotiler(cfg)
	if err != nil {
		panic(err)
	}

	results, err := at.CreateMaps(beach)
	if err != nil {
		panic(err)
	}

	ldr := autotile.NewFileLoader("test/tobs/")
	bin := autotile.NewObjectBin(ldr)
	for _, grp := range tobs {
		err = bin.Load(grp.Group, grp.Chance, grp.Names, grp.All, grp.Any)
		if err != nil {
			panic(err)
		}
	}

	for m := range results {
		at.AddObjects(m, bin)

		err = m.Tilemap.WriteFile(fmt.Sprintf("maptest.01.%d.%d.tmx", m.MapX, m.MapY))
		if err != nil {
			panic(err)
		}
	}

}
