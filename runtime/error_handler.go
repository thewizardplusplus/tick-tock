package runtime

// ErrorHandler ...
//go:generate mockery -name=ErrorHandler -inpkg -case=underscore -testonly
type ErrorHandler interface {
	HandleError(err error)
}

// DefaultErrorHandler ...
type DefaultErrorHandler struct{}

// HandleError ...
func (handler DefaultErrorHandler) HandleError(err error) {
	panic(err)
}
