package tests

import (
	"io"
	"time"

	"github.com/spf13/afero"
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

// Sleeper ...
//go:generate mockery -name=Sleeper -case=underscore
type Sleeper interface {
	Sleep(duration time.Duration)
}

// Exiter ...
//go:generate mockery -name=Exiter -case=underscore
type Exiter interface {
	Exit(code int)
}
