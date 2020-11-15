package runtime

import (
	"io"

	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

//go:generate mockery -name=Writer -inpkg -case=underscore -testonly

// Writer ...
//
// It's used only for mock generating.
//
type Writer interface {
	io.Writer
}

//go:generate mockery -name=Context -inpkg -case=underscore -testonly

// Context ...
//
// It's used only for mock generating.
//
type Context interface {
	context.Context
}

//go:generate mockery -name=ExiterInterface -inpkg -case=underscore -testonly

// ExiterInterface ...
//
// It's used only for mock generating.
//
type ExiterInterface interface {
	Exit(code int)
}
