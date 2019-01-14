package tests

// ...
const (
	UnbufferedInbox = iota
	BufferedInbox
)

// GetNumberAddress ...
func GetNumberAddress(f float64) *float64 {
	return &f
}

// GetStringAddress ...
func GetStringAddress(s string) *string {
	return &s
}
