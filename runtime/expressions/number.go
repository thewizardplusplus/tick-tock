package expressions

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// Number ...
type Number struct {
	value float64
}

// NewNumber ...
func NewNumber(value float64) Number {
	return Number{value}
}

// Evaluate ...
func (expression Number) Evaluate(context context.Context) (result interface{}, err error) {
	return expression.value, nil
}
