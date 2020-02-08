package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestNewFunctionCall(test *testing.T) {
	type fields struct {
		name      string
		arguments []Expression
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
				name: "add",
				arguments: []Expression{
					func() Expression {
						expression := NewSignedExpression("one")
						expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

						return expression
					}(),
					func() Expression {
						expression := NewSignedExpression("two")
						expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(4.2, nil)

						return expression
					}(),
				},
			},
			args: args{
				context: func() context.Context {
					add := func(a float64, b float64) float64 { return a + b }

					context := new(contextmocks.Context)
					context.On("Value", "add").Return(add, true)

					return context
				}(),
			},
			wantResult: 6.5,
			wantErr:    assert.NoError,
		},
		{
			name: "error with an unknown function",
			fields: fields{
				name: "add",
				arguments: []Expression{
					func() Expression {
						expression := NewSignedExpression("one")
						expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

						return expression
					}(),
					func() Expression {
						expression := NewSignedExpression("two")
						expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(4.2, nil)

						return expression
					}(),
				},
			},
			args: args{
				context: func() context.Context {
					context := new(contextmocks.Context)
					context.On("Value", "add").Return(nil, false)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error on function calling (incorrect arguments types)",
			fields: fields{
				name: "add",
				arguments: []Expression{
					func() Expression {
						expression := NewSignedExpression("one")
						expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2, nil)

						return expression
					}(),
					func() Expression {
						expression := NewSignedExpression("two")
						expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(3, nil)

						return expression
					}(),
				},
			},
			args: args{
				context: func() context.Context {
					add := func(a float64, b float64) float64 { panic("not implemented") }

					context := new(contextmocks.Context)
					context.On("Value", "add").Return(add, true)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error on function calling (incorrect results count)",
			fields: fields{
				name: "add",
				arguments: []Expression{
					func() Expression {
						expression := NewSignedExpression("one")
						expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

						return expression
					}(),
					func() Expression {
						expression := NewSignedExpression("two")
						expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(4.2, nil)

						return expression
					}(),
				},
			},
			args: args{
				context: func() context.Context {
					add := func(a float64, b float64) (float64, bool) { return a + b, true }

					context := new(contextmocks.Context)
					context.On("Value", "add").Return(add, true)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			expression := NewFunctionCall(data.fields.name, data.fields.arguments)
			gotResult, gotErr := expression.Evaluate(data.args.context)

			for _, argument := range data.fields.arguments {
				mock.AssertExpectationsForObjects(test, argument)
			}
			mock.AssertExpectationsForObjects(test, data.args.context)
			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}
