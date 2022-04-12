package autotile

import (
	"image"
)

// some kind of special tile intersection that needs extra care
type collisionType string

var (
	collisionWaterfallNS collisionType = "waterfall-ns"
	collisionWaterfallSN collisionType = "waterfall-sn"
	collisionWaterfallEW collisionType = "waterfall-ew"
	collisionWaterfallWE collisionType = "waterfall-we"

	collisionStairsNS collisionType = "stairs-ns"
	collisionStairsSN collisionType = "stairs-sn"
	collisionStairsEW collisionType = "stairs-ew"
	collisionStairsWE collisionType = "stairs-we"
)

func (t collisionType) isStairs() bool {
	switch t {
	case collisionStairsNS, collisionStairsSN, collisionStairsWE, collisionStairsEW:
		return true
	}
	return false
}

//
type collisionHandler struct {
	cols []*collision
}

func (c *collisionHandler) All() []*collision {
	return c.cols
}

func (c *collisionHandler) append(e *Event) {
	for _, col := range c.cols {
		if col.accept(e) {
			return
		}
	}
	newCol := &collision{
		typ:    e.collisionType,
		events: []*Event{e},
		maxX:   e.X,
		minX:   e.X,
		maxY:   e.Y,
		minY:   e.Y,
	}
	c.cols = append(c.cols, newCol)
}

type collision struct {
	typ    collisionType
	events []*Event

	maxX int
	maxY int
	minX int
	minY int
}

func (c *collision) Max() image.Rectangle {
	return image.Rect(c.minX, c.minY, c.maxX, c.maxY)
}

func (c *collision) Events() []*Event {
	return c.events
}

func (c *collision) accept(in *Event) bool {
	if c.typ != in.collisionType {
		return false
	}
	for _, e := range c.events {
		dx := e.X - in.X
		dy := e.Y - in.Y
		if (dx <= 1 && dx >= -1) && (dy <= 1 && dy >= -1) {
			if in.X < c.minX {
				c.minX = in.X
			}
			if in.X > c.maxX {
				c.maxX = in.X
			}
			if in.Y < c.minY {
				c.minY = in.Y
			}
			if in.Y > c.maxY {
				c.maxY = in.Y
			}
			c.events = append(c.events, in)
			return true
		}
	}
	return false
}

func newCollisionHandler() *collisionHandler {
	return &collisionHandler{
		cols: []*collision{},
	}
}
