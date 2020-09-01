package types

// Nil ...
type Nil struct{}

// String ...
func (Nil) String() string {
	return "null"
}

// MarshalJSON ...
func (nilValue Nil) MarshalJSON() (text []byte, err error) {
	return []byte(nilValue.String()), nil
}
