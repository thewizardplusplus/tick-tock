package tests

// ...
const (
	UnbufferedInbox = iota
	BufferedInbox
)

// GetAddress ...
func GetAddress(s string) *string {
	return &s
}
