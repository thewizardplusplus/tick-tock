package types

// Nil ...
type Nil struct{}

// String ...
func (Nil) String() string {
	return "<nil>"
}

// MarshalText ...
func (nilValue Nil) MarshalText() (text []byte, err error) {
	return []byte(nilValue.String()), nil
}
