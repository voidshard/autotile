package estate

import (
	"github.com/voidshard/autotile"
)

// Config holds how an estate should be set up. We consider this in terms of
// What objects should be placed at each level & what sub set(s) of objects
// these contain.
// This is similar to a tree - where the root is the whole estate & we break
// down into leaf nodes (objects to place) and sets (futher nodes in the tree).
// At each set node level we binpack all of the objects into the smallest area
// we can manage & so on down.
type Config struct {
	// Set includes stuff to place. Required.
	Set *Set

	// Seed for RNG purposes. If not set we set one.
	Seed int64

	// Loader to use for tob(s)
	// If not given we set something to read tobs from local disk
	// in the current dir.
	Loader autotile.Loader
}
