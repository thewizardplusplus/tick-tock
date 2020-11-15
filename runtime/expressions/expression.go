package expressions

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

//go:generate mockery -name=Expression -inpkg -case=underscore -testonly

// Expression ...
type Expression interface {
	Evaluate(context context.Context) (result interface{}, err error)
}
