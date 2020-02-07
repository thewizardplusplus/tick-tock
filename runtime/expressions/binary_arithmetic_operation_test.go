package expressions

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/internal/tests"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestNewBinaryArithmeticOperation(test *testing.T) {
	leftOperand := NewSignedExpression("left")
	rightOperand := NewSignedExpression("right")
	handler := func(float64, float64) float64 { panic("not implemented") }
	got := NewBinaryArithmeticOperation(leftOperand, rightOperand, handler)

	mock.AssertExpectationsForObjects(test, leftOperand, rightOperand)
	assert.Equal(test, leftOperand, got.leftOperand)
	assert.Equal(test, rightOperand, got.rightOperand)
	assert.Equal(test, tests.GetReflectionAddress(handler), tests.GetReflectionAddress(got.handler))
}

func TestBinaryArithmeticOperation_Evaluate(test *testing.T) {
	type fields struct {
		leftOperand  Expression
		rightOperand Expression
		handler      BinaryArithmeticOperationHandler
	}
	type args struct {
		context context.Context
	}

	for _, data := range []struct {
		name       string
		fields     fields
		args       args
		wantResult interface{}
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("left")
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

					return expression
				}(),
				rightOperand: func() Expression {
					expression := NewSignedExpression("right")
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(4.2, nil)

					return expression
				}(),
				handler: func(a float64, b float64) float64 { return a + b },
			},
			args:       args{new(contextmocks.Context)},
			wantResult: 6.5,
			wantErr:    assert.NoError,
		},
		{
			name: "error on left operand evaluation",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("left")
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(nil, iotest.ErrTimeout)

					return expression
				}(),
				rightOperand: NewSignedExpression("right"),
				handler:      func(float64, float64) float64 { panic("not implemented") },
			},
			args:       args{new(contextmocks.Context)},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error on right operand evaluation",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("left")
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

					return expression
				}(),
				rightOperand: func() Expression {
					expression := NewSignedExpression("right")
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(nil, iotest.ErrTimeout)

					return expression
				}(),
				handler: func(float64, float64) float64 { panic("not implemented") },
			},
			args:       args{new(contextmocks.Context)},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			expression := BinaryArithmeticOperation{
				leftOperand:  data.fields.leftOperand,
				rightOperand: data.fields.rightOperand,
				handler:      data.fields.handler,
			}
			gotResult, gotErr := expression.Evaluate(data.args.context)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.leftOperand,
				data.fields.rightOperand,
				data.args.context,
			)
			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}
