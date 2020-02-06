package expressions

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// Expression ...
//go:generate mockery -name=Expression -case=underscore
type Expression interface {
	Evaluate(context context.Context) (result interface{}, err error)
}
