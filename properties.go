package autotile

import (
	"github.com/voidshard/tile"
)

const (
	// relative z levels of objects (higher is .. higher)
	zoffsetObject    = 6
	zoffsetWaterfall = 5
	zoffsetCliff     = 4
	zoffsetRoad      = 3
	zoffsetWater     = 2
	zoffsetLand      = 0

	// names of various properties of interest
	pObject = "object"
	pBiome  = "biome"
	pWall   = "wall"
	pWater  = "water"
	pLava   = "lava"
)

var (
	propertiesWater *tile.Properties = nil
	propertiesCliff *tile.Properties = nil
	propertiesWFall *tile.Properties = nil
	propertiesLand  *tile.Properties = nil
	propertiesRoad  *tile.Properties = nil
	propertiesLava  *tile.Properties = nil
)

func init() {
	// in order to be nice we set certain properties on land tiles -- readable
	// from the tileset metadata

	propertiesWater = tile.NewProperties()
	propertiesWater.SetString(pObject, "water")
	propertiesWater.SetString(pBiome, "river")
	propertiesWater.SetBool(pWater, true)

	propertiesLand = tile.NewProperties()
	propertiesLand.SetString(pObject, "land")
	propertiesLand.SetString(pBiome, "land")

	propertiesRoad = tile.NewProperties()
	propertiesRoad.SetString(pObject, "road")
	propertiesRoad.SetString(pBiome, "land")

	propertiesCliff = tile.NewProperties()
	propertiesCliff.SetString(pObject, "cliff-face")
	propertiesCliff.SetString(pBiome, "rock")
	propertiesCliff.SetBool(pWall, true)

	propertiesLava = tile.NewProperties()
	propertiesLava.SetString(pObject, "lava")
	propertiesLava.SetString(pBiome, "volcanic")
	propertiesLava.SetBool(pLava, true)

	propertiesWFall = tile.NewProperties()
	propertiesWFall.SetString(pObject, "waterfall")
	propertiesWFall.SetString(pBiome, "river")
	propertiesWFall.SetBool(pWall, true)
	propertiesWFall.SetBool(pWater, true)
}
