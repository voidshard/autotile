# Autotile

Create tiled maps for an arbitrarily large world space from a simple interface, then add larger objects randomly with simple rules (eg. place trees but only on grass/dirt tiles).

![example beach scene](https://raw.githubusercontent.com/voidshard/autotile/main/assets/beach.01.png)


### Why

Creating maps by hand is labourious and annoying. Especially if you have a few different tilesets for various biomes that need to be placed the same way. A clear river running through an idyllic valley, a murky river through a swamp, a frozen river in the dead of winter - they all have the same tile pieces; centres, left turns, right turns, corners .. like a great jigsaw .. why place these things by hand when a machine can do it (even if it's just to get you started!).


### What

Autotile is intended to place various elements tile by tile on to map(s) from something that implements a simple interface. Including complex things with tiles that dove tail together & can vary in shape; rivers, waterfalls, cliffs, lava, roads / paths. As well as simple land tiles - snow where it is cold, sand on beaches & in deserts, grass, dirt etc. And larger logical objects composed of tiles themselves (eg. trees, rocks, houses & anything really).

For the later, placing larger logical objects, this uses & is designed to be used with my [tile library](https://github.com/voidshard/tile) which includes tooling to disect larger images in to smaller layered maps. For more info & a tool to create compatible 'tobs' please check that [readme](https://github.com/voidshard/tile) or see some of the examples under [tobs](https://github.com/voidshard/autotile/blob/main/test/tobs/tree.large.06.tmx). Tl;dr they're essentially [tiled](https://www.mapeditor.org/) TMX maps with some special properties.


### How

#### Base Maps

Firstly we'd like to create the base maps. By 'base' I mean we want to place all of the landscape type tiles - land tiles like grass, snow, dirt, water, cliffs, lava, waterfalls etc. 

In order to do this we supply something that satisfies the [Outline interface](https://github.com/voidshard/autotile/blob/main/outline.go) (don't panic, it's pretty short). 
```golang
// tell me how large the world map is 
Bounds() image.Rectangle

// tell me what the area is like at world co-ords (x, y)
At(x, y int) *Area
```
An [Area](https://github.com/voidshard/autotile/blob/main/area.go) just relates some basic info like
- height at the given location
- average temperature 
- whether the given point has water (sea, river, swamp), lava, road etc
- a [Land](https://github.com/voidshard/autotile/blob/main/land.go) struct that tells us what tile(s) to place at this location
- optional tags ([]string) to add to tiles sourced from this location

Nb. it's recommended not to re-create a Land struct (or even Area if possible) each time At is called -- we'll be calling this function a lot & having it re-instantiate an object each time is probably more work than we want to do.

Armed with this & some [config](https://github.com/voidshard/autotile/blob/main/config.go) information we can begin creating maps. Checkout the [example](https://github.com/voidshard/autotile/blob/main/test/main.go).


#### Placing Objects

To place tile objects (again, see 

```golang
	ldr := autotile.NewFileLoader("")
	bin := autotile.NewObjectBin(ldr)
  
	bin.Load(
    "trees",  // load a new group called "trees"
    0.4,      // we should place an item from the group "trees" 40% of the time
    []string{"tree.01.tmx", "shrub.01.tmx"}, // here are the trees that the Loader (above) knows how to load
    nil,
    []string{autotile.Dirt, autotile.Grass}, // items from group "trees" can be placed only on Dirt or Grass tiles
  )
  bin.Load(
     "",      // empty string represents the nil group; ie the chance we place nothing at all
     0.6,     // %60 chance we don't place any object on a given tile
     nil, 
     nil, 
     nil,
  )
  
  // where `m` is a MapOutline returned from CreateMaps and `at` is an autotiler from NewAutotiler()
  at.AddObjects(m, bin)

```
Again checkout the [example](https://github.com/voidshard/autotile/blob/main/test/main.go) for more details.

Notice a couple of things
- the Loader here in an interface with one function that loads a TMX map given some string. The most trivial example is FileLoader (where the key is a file path) but of course you can supply your own loader that does whatever
- we can control what tiles the bottom (lowest z-layer) of the object sits on with the last two args `all` and `any` which both take a list of tags ([]string)
- these tags are added at map creation time & as mentioned before the user can stipulate their own (by setting on an Area struct)

Finally we can save out the resulting maps with a simple 
```golang
	m.Tilemap.WriteFile("my-map.tmx")
```
It's uh, recommended to not try to keep thousands of maps in memory at a time & to write the out ASAP.


#### The World

The intention then is to turn a high level world map (depicting rivers , sea, height information, temperature, rivers, sea etc) into an arbitrarily large number of fully tiled maps, each of them representing some (x,y) offset chunk of the world space with fairly minimal work on our part.


### TODO

There's more to come in this space -- I'd like to handle creating interiors, cities & villages, cave systems etc. Feel free to push up PRs, requests, fixes etc. 


### Credits

The image(s) used here are taken from pokemon gen5 tilesets from deviantart and/or spriters resource. Some have been created from existing tiles where original pieces didn't exist. I've added them as an example & a way to visually test code changes.
