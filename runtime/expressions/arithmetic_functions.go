package expressions

import (
	"math"

	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// NewArithmeticNegation ...
func NewArithmeticNegation(operand Expression) FunctionApplying {
	return NewFunctionApplying(
		[]Expression{operand},
		NewArithmeticFunctionHandler(func(context context.Context, arguments []float64) (float64, error) {
			return -arguments[0], nil
		}),
	)
}

// NewMultiplication ...
func NewMultiplication(leftOperand Expression, rightOperand Expression) FunctionApplying {
	return NewFunctionApplying(
		[]Expression{leftOperand, rightOperand},
		NewArithmeticFunctionHandler(func(context context.Context, arguments []float64) (float64, error) {
			return arguments[0] * arguments[1], nil
		}),
	)
}

// NewDivision ...
func NewDivision(leftOperand Expression, rightOperand Expression) FunctionApplying {
	return NewFunctionApplying(
		[]Expression{leftOperand, rightOperand},
		NewArithmeticFunctionHandler(func(context context.Context, arguments []float64) (float64, error) {
			return arguments[0] / arguments[1], nil
		}),
	)
}

// NewModulo ...
func NewModulo(leftOperand Expression, rightOperand Expression) FunctionApplying {
	return NewFunctionApplying(
		[]Expression{leftOperand, rightOperand},
		NewArithmeticFunctionHandler(func(context context.Context, arguments []float64) (float64, error) {
			return math.Mod(arguments[0], arguments[1]), nil
		}),
	)
}

// NewAddition ...
func NewAddition(leftOperand Expression, rightOperand Expression) FunctionApplying {
	return NewFunctionApplying(
		[]Expression{leftOperand, rightOperand},
		NewArithmeticFunctionHandler(func(context context.Context, arguments []float64) (float64, error) {
			return arguments[0] + arguments[1], nil
		}),
	)
}

// NewSubtraction ...
func NewSubtraction(leftOperand Expression, rightOperand Expression) FunctionApplying {
	return NewFunctionApplying(
		[]Expression{leftOperand, rightOperand},
		NewArithmeticFunctionHandler(func(context context.Context, arguments []float64) (float64, error) {
			return arguments[0] - arguments[1], nil
		}),
	)
}
