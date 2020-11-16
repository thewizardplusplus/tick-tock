package interpreter

import (
	"io"

	"github.com/spf13/afero"
	syncutils "github.com/thewizardplusplus/go-sync-utils"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

//go:generate mockery -name=Reader -inpkg -case=underscore -testonly

// Reader ...
//
// It's used only for mock generating.
//
type Reader interface {
	io.Reader
}

//go:generate mockery -name=Context -inpkg -case=underscore -testonly

// Context ...
//
// It's used only for mock generating.
//
type Context interface {
	context.Context
}

//go:generate mockery -name=ErrorHandler -inpkg -case=underscore -testonly

// ErrorHandler ...
//
// It's used only for mock generating.
//
type ErrorHandler interface {
	runtime.ErrorHandler
}

//go:generate mockery -name=WaitGroup -inpkg -case=underscore -testonly

// WaitGroup ...
//
// It's used only for mock generating.
//
type WaitGroup interface {
	syncutils.WaitGroup
}

//go:generate mockery -name=File -inpkg -case=underscore -testonly

// File ...
//
// It's used only for mock generating.
//
type File interface {
	afero.File
}

//go:generate mockery -name=FileSystem -inpkg -case=underscore -testonly

// FileSystem ...
//
// It's used only for mock generating.
//
type FileSystem interface {
	afero.Fs
}
