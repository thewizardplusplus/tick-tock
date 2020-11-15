package expressions

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	expressionsmocks "github.com/thewizardplusplus/tick-tock/runtime/expressions/mocks"
)

type SignedExpression struct {
	*expressionsmocks.Expression

	Sign string
}

func NewSignedExpression(sign string) SignedExpression {
	return SignedExpression{new(expressionsmocks.Expression), sign}
}

func TestNewFunctionCall(test *testing.T) {
	arguments := []Expression{NewSignedExpression("one"), NewSignedExpression("two")}
	got := NewFunctionCall("test", arguments)

	for _, argument := range arguments {
		mock.AssertExpectationsForObjects(test, argument)
	}
	assert.Equal(test, "test", got.name)
	assert.Equal(test, arguments, got.arguments)
}

func TestFunctionCall_Evaluate(test *testing.T) {
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
						expression.On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).Return(2.3, nil)

						return expression
					}(),
					func() Expression {
						expression := NewSignedExpression("two")
						expression.On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).Return(4.2, nil)

						return expression
					}(),
				},
			},
			args: args{
				context: func() context.Context {
					add := func(a float64, b float64) (float64, error) { return a + b, nil }

					context := new(MockContext)
					context.On("Value", "add").Return(add, true)

					return context
				}(),
			},
			wantResult: 6.5,
			wantErr:    assert.NoError,
		},
		{
			name: "success with the empty interface",
			fields: fields{
				name: "add",
				arguments: []Expression{
					func() Expression {
						expression := NewSignedExpression("one")
						expression.On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).Return(2.3, nil)

						return expression
					}(),
					func() Expression {
						expression := NewSignedExpression("two")
						expression.On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).Return(4.2, nil)

						return expression
					}(),
				},
			},
			args: args{
				context: func() context.Context {
					add := func(a interface{}, b interface{}) (float64, error) {
						return a.(float64) + b.(float64), nil
					}

					context := new(MockContext)
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
				name:      "add",
				arguments: []Expression{NewSignedExpression("one"), NewSignedExpression("two")},
			},
			args: args{
				context: func() context.Context {
					context := new(MockContext)
					context.On("Value", "add").Return(nil, false)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with an incorrect function type",
			fields: fields{
				name:      "add",
				arguments: []Expression{NewSignedExpression("one"), NewSignedExpression("two")},
			},
			args: args{
				context: func() context.Context {
					context := new(MockContext)
					context.On("Value", "add").Return("incorrect", true)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with incorrect argument count",
			fields: fields{
				name:      "add",
				arguments: []Expression{NewSignedExpression("one"), NewSignedExpression("two")},
			},
			args: args{
				context: func() context.Context {
					add := func(a float64, b float64, c float64) (float64, error) { return a + b + c, nil }

					context := new(MockContext)
					context.On("Value", "add").Return(add, true)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with incorrect result count",
			fields: fields{
				name:      "add",
				arguments: []Expression{NewSignedExpression("one"), NewSignedExpression("two")},
			},
			args: args{
				context: func() context.Context {
					add := func(a float64, b float64) float64 { return a + b }

					context := new(MockContext)
					context.On("Value", "add").Return(add, true)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with an incorrect result type",
			fields: fields{
				name:      "add",
				arguments: []Expression{NewSignedExpression("one"), NewSignedExpression("two")},
			},
			args: args{
				context: func() context.Context {
					add := func(a float64, b float64) (float64, bool) { return a + b, true }

					context := new(MockContext)
					context.On("Value", "add").Return(add, true)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with argument evaluation",
			fields: fields{
				name: "add",
				arguments: []Expression{
					func() Expression {
						expression := NewSignedExpression("one")
						expression.
							On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).
							Return(nil, iotest.ErrTimeout)

						return expression
					}(),
					NewSignedExpression("two"),
				},
			},
			args: args{
				context: func() context.Context {
					add := func(a float64, b float64) (float64, error) { return a + b, nil }

					context := new(MockContext)
					context.On("Value", "add").Return(add, true)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with an incorrect argument type",
			fields: fields{
				name: "add",
				arguments: []Expression{
					func() Expression {
						expression := NewSignedExpression("one")
						expression.On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).Return(2, nil)

						return expression
					}(),
					NewSignedExpression("two"),
				},
			},
			args: args{
				context: func() context.Context {
					add := func(a float64, b float64) (float64, error) { return a + b, nil }

					context := new(MockContext)
					context.On("Value", "add").Return(add, true)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with function calling",
			fields: fields{
				name: "add",
				arguments: []Expression{
					func() Expression {
						expression := NewSignedExpression("one")
						expression.On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).Return(2.3, nil)

						return expression
					}(),
					func() Expression {
						expression := NewSignedExpression("two")
						expression.On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).Return(4.2, nil)

						return expression
					}(),
				},
			},
			args: args{
				context: func() context.Context {
					add := func(a float64, b float64) (float64, error) { return 0, iotest.ErrTimeout }

					context := new(MockContext)
					context.On("Value", "add").Return(add, true)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			expression := FunctionCall{
				name:      data.fields.name,
				arguments: data.fields.arguments,
			}
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
