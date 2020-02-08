package expressions

import (
	"math"
)

// NewMultiplication ...
func NewMultiplication(leftOperand Expression, rightOperand Expression) FunctionApplying {
	return NewFunctionApplying(
		[]Expression{leftOperand, rightOperand},
		NewArithmeticFunctionHandler(func(arguments []float64) (float64, error) {
			return arguments[0] * arguments[1], nil
		}),
	)
}

// NewDivision ...
func NewDivision(leftOperand Expression, rightOperand Expression) FunctionApplying {
	return NewFunctionApplying(
		[]Expression{leftOperand, rightOperand},
		NewArithmeticFunctionHandler(func(arguments []float64) (float64, error) {
			return arguments[0] / arguments[1], nil
		}),
	)
}

// NewModulo ...
func NewModulo(leftOperand Expression, rightOperand Expression) FunctionApplying {
	return NewFunctionApplying(
		[]Expression{leftOperand, rightOperand},
		NewArithmeticFunctionHandler(func(arguments []float64) (float64, error) {
			return math.Mod(arguments[0], arguments[1]), nil
		}),
	)
}
