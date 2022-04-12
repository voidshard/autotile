package estate

// Object represents an object we can place in an estate
type Object struct {
	// Identifier for 'tob' tmx object file, used with Loader interface
	Tob string

	// Max number of items
	// If max != min a random number between the two will be chosen
	MaxCount int

	// Min number of items
	// If max != min a random number between the two will be chosen.
	// If 0 or a negative is chosen we disown this object all together
	// (ie, no fence, base or anything from this will be set)
	MinCount int
}

// NewObject makes a new obj with min/max of 1
func NewObject(tob string) *Object {
	return &Object{Tob: tob, MinCount: 1, MaxCount: 1}
}

// SetMax sets the max count & returns the obj as a convience func.
func (o *Object) SetMax(i int) *Object {
	o.MaxCount = i
	return o
}

// SetMin sets the min count & returns the obj as a convience func.
func (o *Object) SetMin(i int) *Object {
	o.MinCount = i
	return o
}
