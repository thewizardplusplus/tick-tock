package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestArithmeticFunctions(test *testing.T) {
	type args struct {
		context context.Context
	}

	for _, data := range []struct {
		name       string
		expression FunctionApplying
		args       args
		wantResult interface{}
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "arithmetic negation",
			expression: func() FunctionApplying {
				operand := NewSignedExpression("left")
				operand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.0, nil)

				return NewArithmeticNegation(operand)
			}(),
			args:       args{new(contextmocks.Context)},
			wantResult: -2.0,
			wantErr:    assert.NoError,
		},
		{
			name: "multiplication",
			expression: func() FunctionApplying {
				leftOperand := NewSignedExpression("left")
				leftOperand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.0, nil)

				rightOperand := NewSignedExpression("right")
				rightOperand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(3.0, nil)

				return NewMultiplication(leftOperand, rightOperand)
			}(),
			args:       args{new(contextmocks.Context)},
			wantResult: 6.0,
			wantErr:    assert.NoError,
		},
		{
			name: "division",
			expression: func() FunctionApplying {
				leftOperand := NewSignedExpression("left")
				leftOperand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(10.0, nil)

				rightOperand := NewSignedExpression("right")
				rightOperand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.0, nil)

				return NewDivision(leftOperand, rightOperand)
			}(),
			args:       args{new(contextmocks.Context)},
			wantResult: 5.0,
			wantErr:    assert.NoError,
		},
		{
			name: "modulo",
			expression: func() FunctionApplying {
				leftOperand := NewSignedExpression("left")
				leftOperand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(5.0, nil)

				rightOperand := NewSignedExpression("right")
				rightOperand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.0, nil)

				return NewModulo(leftOperand, rightOperand)
			}(),
			args:       args{new(contextmocks.Context)},
			wantResult: 1.0,
			wantErr:    assert.NoError,
		},
		{
			name: "addition",
			expression: func() FunctionApplying {
				leftOperand := NewSignedExpression("left")
				leftOperand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.0, nil)

				rightOperand := NewSignedExpression("right")
				rightOperand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(3.0, nil)

				return NewAddition(leftOperand, rightOperand)
			}(),
			args:       args{new(contextmocks.Context)},
			wantResult: 5.0,
			wantErr:    assert.NoError,
		},
		{
			name: "subtraction",
			expression: func() FunctionApplying {
				leftOperand := NewSignedExpression("left")
				leftOperand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(5.0, nil)

				rightOperand := NewSignedExpression("right")
				rightOperand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.0, nil)

				return NewSubtraction(leftOperand, rightOperand)
			}(),
			args:       args{new(contextmocks.Context)},
			wantResult: 3.0,
			wantErr:    assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotResult, gotErr := data.expression.Evaluate(data.args.context)

			for _, argument := range data.expression.arguments {
				mock.AssertExpectationsForObjects(test, argument)
			}
			mock.AssertExpectationsForObjects(test, data.args.context)
			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}
