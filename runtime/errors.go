package runtime

import "github.com/pkg/errors"

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
type DefaultErrorHandler struct{}

// HandleError ...
func (handler DefaultErrorHandler) HandleError(err error) {
	panic(err)
}
