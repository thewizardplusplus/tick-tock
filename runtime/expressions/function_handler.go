package expressions

import (
	"github.com/pkg/errors"
)

// FunctionHandler ...
type FunctionHandler func([]interface{}) (interface{}, error)

// ArithmeticFunctionHandler ...
type ArithmeticFunctionHandler func([]float64) (float64, error)

// NewArithmeticFunctionHandler ...
func NewArithmeticFunctionHandler(handler ArithmeticFunctionHandler) FunctionHandler {
	return func(arguments []interface{}) (interface{}, error) {
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

		return handler(values)
	}
}
