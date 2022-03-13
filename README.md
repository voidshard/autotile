# Autotile

Create tiled maps for an arbitrarily large world space from a simple interface, then add larger objects randomly with simple rules (eg. place trees but only on grass/dirt tiles).

![example beach scene](https://raw.githubusercontent.com/voidshard/autotile/main/assets/beach.01.png)


### Why

Creating maps by hand is labourious and annoying. Especially if you have a few different tilesets for various biomes that need to be placed the same way. A clear river running through an idyllic valley, a murky river through a swamp, a frozen river in the dead of winter, even lava flowing down a volcano - they all have the same tile pieces; centres, left turns, right turns, corners like a great jigsaw. The problem only gets worse when you consider all the other multi-tile objects you might want to place: trees, houses, castles ... oof. Ideally we'd tell something "here is a sketch of the map I want, here are my objects, do the tiling for me!" .. and that's what Autotile is supposed to do.


### What

Autotile is intended to place various elements tile by tile on to map(s) from something that implements a simple interface. Including complex things with tiles that dove tail together & can vary in shape; rivers, waterfalls, cliffs, lava, roads / paths. As well as simple land tiles - snow where it is cold, sand on beaches & in deserts, grass, dirt etc. And larger logical objects composed of tiles themselves (eg. trees, rocks, houses & anything really).

For the later, placing larger logical objects, this uses & is designed to be used with my [tile library](https://github.com/voidshard/tile) which includes tooling to disect larger images in to smaller layered maps. For more info & a tool to create compatible 'tobs' please check that [readme](https://github.com/voidshard/tile) or see some of the examples under [tobs](https://github.com/voidshard/autotile/blob/main/test/tobs/tree.large.06.tmx). Tl;dr they're essentially [tiled](https://www.mapeditor.org/) TMX maps with some special properties.


### How

#### Base Maps

Firstly we'd like apply base or landscape type tiles to our map - grass, snow, dirt, water, cliffs, lava, waterfalls etc.

In order to do this we supply something that satisfies the [Outline interface](https://github.com/voidshard/autotile/blob/main/interface.go) 
```golang
// tell me what the area is like at world co-ords (x, y)
LandAt(x, y int) LandData
```
Returning [LandData](https://github.com/voidshard/autotile/blob/main/interface.go) for the given location to answer some basic questions
- height 
- average temperature 
- whether the given point has water (sea, river, swamp), lava, is a road etc
- a [LandTiles](https://github.com/voidshard/autotile/blob/main/landtiles.go) struct that tells us what tile(s) we can place at this location
- optional tags ([]string) associated with this (x,y) (some [tags](https://github.com/voidshard/autotile/blob/main/tags.go) are set by the tiler but this allows user defined tags)

Armed with this & some [config](https://github.com/voidshard/autotile/blob/main/config.go) information we can begin tiling maps. Checkout the [example](https://github.com/voidshard/autotile/blob/main/test/main.go).

```golang
  // prep a 30x30 tile map, where tiles are 32x32 pixels
  tmap := tile.New(&tile.Config{ // from github.com/voidshard/tile
    TileWidth:  32,
    TileHeight: 32,
    MapWidth:   30,
    MapHeight:  30, 
  })

  // outline some config information about where various types of 
  // terrain occur. The numbers here depend on what your Outline returns.
  // Eg. if when you return Height() you mean 'cm' then saying a mountain 
  // is anything over 240 is .. odd.
  var cfg = &autotile.Config{
    BeachWidth:        2,
    VegetationMaxTemp: 45,
    VegetationMinTemp: -5,
    MountainLevel:     240,
    CliffLevel:        170,
  }

  // prep the autotiler
  at, err := autotile.NewAutotiler(cfg)
  if err != nil {
    panic(err)
  }

  // lay the base tiles (where 'beach' here implements Outline)
  err = at.SetLand(beach, beach.Bounds(), tmap)
  if err != nil {
    panic(err)
  }
```


#### Placing Objects

Now that we have our base tiles down, we can go ahead and place static objects (trees, houses ..). For this the autotiler provides the SetObjects() function which takes 
- our Outline interface
- a image.Rectangle "bounds" of where to place objects
- a tile map to write to
- an 'ObjectBin' that can choose what object(s) to place where

```golang
  // loader that reads .tmx objects from disk from current dir
  ldr := autotile.NewFileLoader("")

  bin := autotile.NewBin(at, beach, 987654321, ldr)
  
  bin.Load(
    "trees",  // load a new group called "trees"
    &autotile.LoadConfig{
      Chance: 0.4,      // we should place an item from the group "trees" 40% of the time
      Objects: []string{"tree.01.tmx", "shrub.01.tmx"}, // here are the trees that the Loader (above) knows how to load
      TagsAny: []string{autotile.Dirt, autotile.Grass}, // items from group "trees" can be placed only on Dirt or Grass tiles
      Distribution: autotile.RandomDistribution, // layout trees randomly
    },
  )
  
  bin.Load(
     "",      // empty string represents the nil group; ie the chance we place nothing at all
     &autotile.LoadConfig{
        Chance: 0.6,    // %60 chance we don't place any object on a given tile
     }
  )
  
  // place objects loaded in `bin` on `tmap` within the given bounds
  err = at.SetObjects(beach, beach.Bounds(), tmap, bin)
```
Again checkout the [example](https://github.com/voidshard/autotile/blob/main/test/main.go) for more details.

Finally we can save out the resulting maps with a simple 
```golang
	tmap.WriteFile("my-map.tmx")
```

#### The World

The intention then is to turn a high level world map (depicting rivers, sea, height information, temperature, lava, swamps etc) into an arbitrarily large number of fully tiled maps, each of them representing some (x,y) offset chunk of the world space with fairly minimal work on our part. For a simple example toy lib for this I have some [trivial worldgen code](https://github.com/voidshard/cartographer/blob/master/pkg/landscape/perlinworld.go).


### Notes

Complex objects like cliffs, water etc take certain numbers of tiles to place. Eg. rivers need to place a left bank, right bank and (hopefully) a centre, so you need 2-3 tiles of side by side 'water' tiles to set nicely. For the same reason you'll want to make sure that rivers keep a valid width when flowing diagonally.

Waterfalls are currently only supported North to South (flowing toward the viewer) or South to North (flowing away). If you have example tiles of a waterfall flowing East-West or West-East let me know & I can add.

If you want to tile a whole world (who doesn't) then you probably need to go map by map, implement the tile.Tileable interface with something clever w.r.t. memory management or try the trivial [infinite map](https://github.com/voidshard/tile/blob/master/infinite.go) (which keeps the "map" in a tempfile on disk so it can be arbitrarily large).

It's recommended not to do too much work when LandAt is called, we'll be calling it a lot & it's performance drastically alters map tiling time(s).

Some things to note on object placement
- the Loader here is an interface with one function that loads a TMX map given some string. The most trivial example is FileLoader (where the key is a file path) but of course you can supply your own loader that does whatever
- the ObjectBin here is another interface with one function that chooses an object (TMX) to place given a proposed destination. The Bin is fairly simple, you can of course supply your own
- we can control what tiles the bottom (lowest z-layer) of an object sits on with tags `TagsAll` and `TagsAny` which both take a list of tags ([]string)
- default tags (seen in examples) are added at map creation time (see [tags.go](https://github.com/voidshard/autotile/blob/main/tags.go)) but the user can stipulate their own additional tags and use these to place objects.
- the provided Bin implementation will not place an object if it would overwrite existing tiles, for this reason smaller objects are easier to place & you may need to adjust probabilities accordingly
- we can supply `Distribution` to indicate how we want random values chosen for a given group. Currently we support `RandomDistribution` & `PerlinDistribution`


### TODO

API might change around for a bit while I'm adding features / organising things.
- 2022-03-13 API has indeed changed to accept the tile.Tileable interface, allowing us to support the new InfiniteMap in the tile lib
- Config struct changed to remove WorldParams as it's own struct

There's more to come in this space -- I'd like to handle creating interiors, cities & villages, cave systems etc. Feel free to push up PRs, requests, fixes etc. 


### Credits

The image(s) used here are taken from pokemon gen5 tilesets from deviantart and/or spriters resource. Some have been created from existing tiles where original pieces didn't exist. I've added them as an example & a way to visually test code changes.
