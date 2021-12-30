package autotile

// Heading represents some compass direction.
type Heading int

const (
	// We start at North and move clockwise.
	North     Heading = 0
	NorthEast Heading = 1
	East      Heading = 2
	SouthEast Heading = 3
	South     Heading = 4
	SouthWest Heading = 5
	West      Heading = 6
	NorthWest Heading = 7
)
