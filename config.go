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
	// Each region will get it's own Seed based on this + some math operations
	// on the min/max x/y values given (deterministically) .. this means that
	// the same function call with the same input Outline for a given rectangle
	// ("region") should choose the same tiles (assuming the Outline is returning
	// the same data each time).
	Seed int64

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

	// BeachWidth is how many tiles sand should travel up from water tiles.
	// Higher values means wider beaches / more sand.
	// Minimum of 0
	BeachWidth int

	// TransitionWidth indicates how many `units` we take to transition from one
	// ground type to another.
	// Ie. we start using sand transitions at VegetationMaxTemp and move
	// to full sand at VegetationMaxTemp+TransitionWidth
	// Minimum of 0
	TransitionWidth int

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

// Validate that the config is correct
func (c *Config) Validate() error {
	if c.Seed == 0 {
		c.Seed = time.Now().UnixNano()
	}

	if c.BeachWidth < 0 {
		c.BeachWidth = 0
	}
	if c.TransitionWidth < 0 {
		c.TransitionWidth = 0
	}
	if c.VegetationMaxTemp <= c.VegetationMinTemp {
		return fmt.Errorf("%w: vegetation max temp should be greater than min temp", ErrInvalidValue)
	}

	if c.ZOffsetLand < 0 {
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

	return nil
}
