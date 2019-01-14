package tests

import "io"

// Writer ...
//go:generate mockery -name=Writer -case=underscore
type Writer interface {
	io.Writer
}

// GetAddress ...
func GetAddress(s string) *string {
	return &s
}
