package runtime

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// ErrReturn ...
var ErrReturn = errors.New("return")

//go:generate mockery --name=ErrorHandler --inpackage --case=underscore --testonly

// ErrorHandler ...
type ErrorHandler interface {
	HandleError(err error)
}

// Exiter ...
type Exiter func(code int)

// DefaultErrorHandler ...
type DefaultErrorHandler struct {
	writer io.Writer
	exiter Exiter
}

// NewDefaultErrorHandler ...
func NewDefaultErrorHandler(writer io.Writer, exiter Exiter) DefaultErrorHandler {
	return DefaultErrorHandler{writer, exiter}
}

// HandleError ...
func (handler DefaultErrorHandler) HandleError(err error) {
	fmt.Fprintf(handler.writer, "error: %s\n", err) // nolint: errcheck, gosec
	handler.exiter(1)
}

func newUnknownStateError(state context.State) error {
	return errors.Errorf("unknown state %s", state.Name)
}
