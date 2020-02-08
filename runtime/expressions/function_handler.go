package expressions

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// FunctionHandler ...
type FunctionHandler func(context context.Context, arguments []interface{}) (interface{}, error)

// ArithmeticFunctionHandler ...
type ArithmeticFunctionHandler func(context context.Context, arguments []float64) (float64, error)

// NewArithmeticFunctionHandler ...
func NewArithmeticFunctionHandler(handler ArithmeticFunctionHandler) FunctionHandler {
	return func(context context.Context, arguments []interface{}) (interface{}, error) {
		var values []float64
		for index, argument := range arguments {
			value, ok := argument.(float64)
			if !ok {
				return nil, errors.Errorf(
					"incorrect type of the argument #%d (%T instead float64)",
					index,
					argument,
				)
			}

			values = append(values, value)
		}

		result, err := handler(context, values)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to call the arithmetic function handler")
		}

		return result, nil
	}
}
