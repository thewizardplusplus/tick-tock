package tests

import (
	"io"

	"github.com/spf13/afero"
)

// ...
const (
	UnbufferedInbox = iota
	BufferedInbox
)

// Writer ...
//go:generate mockery -name=Writer -case=underscore
type Writer interface {
	io.Writer
}

// FileSystem ...
//go:generate mockery -name=FileSystem -case=underscore
type FileSystem interface {
	afero.Fs
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
