package runtime

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// ErrUserExit ...
var ErrUserExit = errors.New("user exit")

func newUnknownStateError(state string) error {
	return errors.Errorf("unknown state %s", state)
}

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
		handler.writer.Write([]byte(fmt.Sprintf("error: %s", err))) // nolint: gosec
		code = 1
	}

	handler.exiter(code)
}
