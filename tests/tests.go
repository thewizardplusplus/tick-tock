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

// Reader ...
//go:generate mockery -name=Reader -case=underscore
type Reader interface {
	io.Reader
}

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

// File ...
//go:generate mockery -name=File -case=underscore
type File interface {
	afero.File
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
