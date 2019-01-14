package runtime

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// ErrUserExit ...
var ErrUserExit = errors.New("user exit")

// ErrorHandler ...
//go:generate mockery -name=ErrorHandler -case=underscore
type ErrorHandler interface {
	HandleError(err error)
}

// DefaultErrorHandler ...
type DefaultErrorHandler struct {
	writer io.Writer
	exiter func(code int)
}

// NewDefaultErrorHandler ...
func NewDefaultErrorHandler(writer io.Writer, exiter func(code int)) DefaultErrorHandler {
	return DefaultErrorHandler{writer, exiter}
}

// HandleError ...
func (handler DefaultErrorHandler) HandleError(err error) {
	var code int
	if errors.Cause(err) != ErrUserExit {
		handler.writer.Write([]byte(fmt.Sprintf("error: %s\n", err))) // nolint: gosec, errcheck
		code = 1
	}

	handler.exiter(code)
}

func newUnknownStateError(state string) error {
	return errors.Errorf("unknown state %s", state)
}
