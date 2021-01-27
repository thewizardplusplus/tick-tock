package options

import (
	"io"
)

//go:generate mockery --name=Writer --inpackage --case=underscore --testonly

// Writer ...
//
// It's used only for mock generating.
//
type Writer interface {
	io.Writer
}

//go:generate mockery --name=ExiterInterface --inpackage --case=underscore --testonly

// ExiterInterface ...
//
// It's used only for mock generating.
//
type ExiterInterface interface {
	Exit(code int)
}
