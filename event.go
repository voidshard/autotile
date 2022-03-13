package autotile

import (
	"github.com/voidshard/tile"
)

//
type Event struct {
	X int
	Y int
	Z int

	Src        string
	Properties *tile.Properties

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
