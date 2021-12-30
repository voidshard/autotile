package autotile

import (
	"fmt"
	"time"
)

var (
	// ErrMissingRequiredValue implies a config / setting should be set
	// that is not
	ErrMissingRequiredValue = "missing required value"

	// ErrInvalidValue means some input is nonsense
	ErrInvalidValue = "input value invalid"
)

// Config dictates various top level concerns with our rendering / creation process
type Config struct {
	// Seed is used for RNG. If zero a random value will be used.
	// Each MapOutline is given the seed mutated deterministically by the map
	// offset (ie. the world map x,y values) for it's own RNG.
	Seed int64

	// step size is how many pixels should be used at a time to make
	// each map. We then enlarge each `StepSize` pixel area into `MapSize` in tiles.
	// Ie. 100 would turn a 1k by 1k world map into 10x10 smaller
	// maps, where each represented an enlarged area of 100x100.
	//
	// Nb. because we need 3 tiles to depict rivers & waterfalls we require
	// that each world map level pixel be expanded into at least 3 area level tiles.
	// This implies that the ratio of map size to step size should be at least 3.
	// In other words "MapSize >= StepSize * 3" should be true.
	// Otherwise .. we can't nicely tile rivers, cliffs, lava, waterfalls etc.
	//
	// A value of 0 would mean 'use the whole input map'
	StepSize int

	// MapSize in tiles. Dictates how many tiles wide & high each
	// map should be. This is expected to be at least 3x `StepSize`
	// (value is enforced).
	//
	// Nb. we don't enforce a max output size for the map -- you can try to make
	// a 1Mx1M tile map if you want. Probably your computer will melt though.
	MapSize int

	// TileSize is the height & width of each tile, in pixels.
	// Defaults to 32.
	TileSize int

	// Params around limits, world settings, biome transitions
	Params *WorldParams

	// Worker routines (number of maps to build simultaneously).
	Routines int

	// Layer at which base land tiles are set.
	// Set to default value if not set. In general you shouldn't need to set this.
	ZOffsetLand int

	// Layer at which water & lava tiles are placed
	// Set to default value if not set. In general you shouldn't need to set this.
	ZOffsetWater int

	// Layer at which road tiles are placed
	// Set to default value if not set. In general you shouldn't need to set this.
	ZOffsetRoad int

	// Layer at this cliff tiles are placed
	// Set to default value if not set. In general you shouldn't need to set this.
	ZOffsetCliff int

	// Layer at which waterfall tiles are placed
	// Set to default value if not set. In general you shouldn't need to set this.
	ZOffsetWaterfall int

	// Layer at which objects are placed (ie, from an objectbin)
	// Set to default value if not set. If you need to set this then you probably want
	// it to be higher than all of the other offsets ..
	ZOffsetObject int
}

// WorldParams lays out some limits or general outlines for the world.
type WorldParams struct {
	// BeachWidth is how many tiles sand should travel up from sea tiles.
	// Higher values means wider beaches / more sand.
	// Minimum of 0
	BeachWidth int

	// VegetationMaxTemp is the highest the temperature can be before we decide
	// it's too hot for there to be vegetation.
	// Temp in degrees C
	VegetationMaxTemp int

	// VegetationMinTemp is how low the temperature can be before we decide it's
	// too cold for there to be vegetation.
	// Temp in degrees C
	VegetationMinTemp int

	// Height above which we consider terrain mountainous (limited to no vegetation)
	// and generally barren / rocky.
	// Height value 0-255
	MountainLevel int

	// CliffLevel is the height at which we consider placing cliffs & waterfalls.
	// Height value 0-255
	CliffLevel int

	// Temperature in degrees below which we render snow / ice
	// Temp in degrees C
	SnowLevel int
}

// validate correct WorldParams
func (w *WorldParams) validate() error {
	if w.BeachWidth < 0 {
		w.BeachWidth = 0
	}
	if w.VegetationMaxTemp <= w.VegetationMinTemp {
		return fmt.Errorf("%w: vegetation max temp should be greater than min temp", ErrInvalidValue)
	}
	return nil
}

// Validate that the config is correct
func (c *Config) Validate() error {
	if c.TileSize < 1 {
		c.TileSize = 32
	}

	if c.Routines < 1 {
		c.Routines = 1
	}

	if c.Seed == 0 {
		c.Seed = time.Now().UnixNano()
	}

	if c.StepSize < 0 {
		c.StepSize = 0
		c.MapSize = 0
	} else if c.MapSize < c.StepSize*3 {
		return fmt.Errorf("%w: MapSize should be >= StepSize * 3", ErrInvalidValue)
	}

	if c.Params == nil {
		c.Params = &WorldParams{
			BeachWidth:        3,
			VegetationMaxTemp: 50,
			VegetationMinTemp: -20,
		}
	}

	if c.ZOffsetLand <= 0 {
		c.ZOffsetLand = zoffsetLand
	}
	if c.ZOffsetWater <= 0 {
		c.ZOffsetWater = zoffsetWater
	}
	if c.ZOffsetRoad <= 0 {
		c.ZOffsetRoad = zoffsetRoad
	}
	if c.ZOffsetCliff <= 0 {
		c.ZOffsetCliff = zoffsetCliff
	}
	if c.ZOffsetWaterfall <= 0 {
		c.ZOffsetWaterfall = zoffsetWaterfall
	}
	if c.ZOffsetObject <= 0 {
		c.ZOffsetObject = zoffsetObject
	}

	return c.Params.validate()
}
