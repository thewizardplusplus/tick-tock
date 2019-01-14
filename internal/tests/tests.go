package tests

// ...
const (
	UnbufferedInbox = iota
	BufferedInbox
)

// GetStringAddress ...
func GetStringAddress(s string) *string {
	return &s
}
