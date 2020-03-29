package expressions

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// String ...
type String struct {
	value *types.Pair
}

// NewString ...
func NewString(value string) String {
	pair := types.NewPairFromText(value)
	return String{pair}
}

// Evaluate ...
func (expression String) Evaluate(context context.Context) (result interface{}, err error) {
	return expression.value, nil
}
