package autotile

import (
	"github.com/voidshard/tile"
)

// Event reports that the autotiler has made a decision & what that decision is.
// We report the location (x,y,z) and either of
// - the id of the object placed
// - the source (Src) & properties set
type Event struct {
	X int
	Y int
	Z int

	// if a tile is set
	Src        string
	Properties *tile.Properties

	// if an object is placed
	ObjectID string
}

//
func newEvent(x, y, z int, src string, props *tile.Properties) *Event {
	return &Event{
		X:          x,
		Y:          y,
		Z:          z,
		Src:        src,
		Properties: props,
	}
}

func newObjEvent(x, y, z int, id string) *Event {
	return &Event{X: x, Y: y, Z: z, ObjectID: id}
}
