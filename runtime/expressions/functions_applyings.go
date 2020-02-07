package expressions

// NewMultiplication ...
func NewMultiplication(leftOperand Expression, rightOperand Expression) FunctionApplying {
	return NewFunctionApplying(
		[]Expression{leftOperand, rightOperand},
		NewArithmeticFunctionHandler(func(arguments []float64) (float64, error) {
			return arguments[0] * arguments[1], nil
		}),
	)
}
