package expressions

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions/mocks"
)

func TestNewArithmeticNegation(test *testing.T) {
	operand := new(mocks.Expression)
	got := NewArithmeticNegation(operand)

	mock.AssertExpectationsForObjects(test, operand)
	assert.Equal(test, operand, got.operand)
}

func TestArithmeticNegation_Evaluate(test *testing.T) {
	type fields struct {
		operand Expression
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
				operand: func() Expression {
					expression := new(mocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

					return expression
				}(),
			},
			args:       args{new(contextmocks.Context)},
			wantResult: -2.3,
			wantErr:    assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				operand: func() Expression {
					expression := new(mocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(nil, iotest.ErrTimeout)

					return expression
				}(),
			},
			args:       args{new(contextmocks.Context)},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			expression := ArithmeticNegation{
				operand: data.fields.operand,
			}
			gotResult, gotErr := expression.Evaluate(data.args.context)

			mock.AssertExpectationsForObjects(test, data.fields.operand, data.args.context)
			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}

func Test_evaluateFloat64Operand(test *testing.T) {
	type args struct {
		context context.Context
		operand Expression
	}

	for _, data := range []struct {
		name      string
		args      args
		wantValue float64
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				context: new(contextmocks.Context),
				operand: func() Expression {
					expression := new(mocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

					return expression
				}(),
			},
			wantValue: 2.3,
			wantErr:   assert.NoError,
		},
		{
			name: "error on operand evaluation",
			args: args{
				context: new(contextmocks.Context),
				operand: func() Expression {
					expression := new(mocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(nil, iotest.ErrTimeout)

					return expression
				}(),
			},
			wantValue: 0,
			wantErr:   assert.Error,
		},
		{
			name: "error on operand typecasting",
			args: args{
				context: new(contextmocks.Context),
				operand: func() Expression {
					expression := new(mocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2, nil)

					return expression
				}(),
			},
			wantValue: 0,
			wantErr:   assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotValue, gotErr := evaluateFloat64Operand(data.args.context, data.args.operand)

			mock.AssertExpectationsForObjects(test, data.args.context, data.args.operand)
			assert.Equal(test, data.wantValue, gotValue)
			data.wantErr(test, gotErr)
		})
	}
}
