package types

// Nil ...
type Nil struct{}

// MarshalJSON ...
func (Nil) MarshalJSON() (text []byte, err error) {
	return []byte("null"), nil
}
