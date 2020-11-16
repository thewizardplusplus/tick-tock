package options

import (
	"io"
)

//go:generate mockery -name=Writer -inpkg -case=underscore -testonly

// Writer ...
//
// It's used only for mock generating.
//
type Writer interface {
	io.Writer
}
