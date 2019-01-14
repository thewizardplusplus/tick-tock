package tests

import "io"

// Writer ...
//go:generate mockery -name=Writer -case=underscore
type Writer interface {
	io.Writer
}

// Exiter ...
//go:generate mockery -name=Exiter -case=underscore
type Exiter interface {
	Exit(code int)
}

// GetAddress ...
func GetAddress(s string) *string {
	return &s
}
