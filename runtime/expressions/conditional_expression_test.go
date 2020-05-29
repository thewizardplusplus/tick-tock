package expressions

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	runtimemocks "github.com/thewizardplusplus/tick-tock/runtime/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

type SignedCommand struct {
	*runtimemocks.Command

	Sign string
}

func NewSignedCommand(sign string) SignedCommand {
	return SignedCommand{new(runtimemocks.Command), sign}
}

func TestNewConditionalExpression(test *testing.T) {
	conditionalCases := []ConditionalCase{
		{NewSignedExpression("one-condition"), NewSignedCommand("one-command")},
		{NewSignedExpression("two-condition"), NewSignedCommand("two-command")},
	}
	got := NewConditionalExpression(conditionalCases)

	checkConditionalCases(test, conditionalCases)
	assert.Equal(test, conditionalCases, got.conditionalCases)
}

func TestConditionalExpression_Evaluate(test *testing.T) {
	type fields struct {
		conditionalCases []ConditionalCase
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
			name: "success/without conditions",
			fields: fields{
				conditionalCases: nil,
			},
			args: args{
				context: new(contextmocks.Context),
			},
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name: "success/with conditions/with match",
			fields: fields{
				conditionalCases: []ConditionalCase{
					{
						Condition: func() Expression {
							expression := NewSignedExpression("one-condition")
							expression.
								On("Evaluate", mock.AnythingOfType("*mocks.Context")).
								Return(types.False, nil)

							return expression
						}(),
						Command: NewSignedCommand("one-command"),
					},
					{
						Condition: func() Expression {
							expression := NewSignedExpression("two-condition")
							expression.
								On("Evaluate", mock.AnythingOfType("*mocks.Context")).
								Return(types.True, nil)

							return expression
						}(),
						Command: func() runtime.Command {
							command := NewSignedCommand("two-command")
							command.On("Run", mock.AnythingOfType("*mocks.Context")).Return(23.0, nil)

							return command
						}(),
					},
					{
						Condition: NewSignedExpression("three-condition"),
						Command:   NewSignedCommand("three-command"),
					},
				},
			},
			args: args{
				context: func() context.Context {
					context := new(contextmocks.Context)
					context.On("Copy").Return(context)

					return context
				}(),
			},
			wantResult: 23.0,
			wantErr:    assert.NoError,
		},
		{
			name: "success/with conditions/without match",
			fields: fields{
				conditionalCases: []ConditionalCase{
					{
						Condition: func() Expression {
							expression := NewSignedExpression("one-condition")
							expression.
								On("Evaluate", mock.AnythingOfType("*mocks.Context")).
								Return(types.False, nil)

							return expression
						}(),
						Command: NewSignedCommand("one-command"),
					},
					{
						Condition: func() Expression {
							expression := NewSignedExpression("two-condition")
							expression.
								On("Evaluate", mock.AnythingOfType("*mocks.Context")).
								Return(types.False, nil)

							return expression
						}(),
						Command: NewSignedCommand("two-command"),
					},
					{
						Condition: func() Expression {
							expression := NewSignedExpression("three-condition")
							expression.
								On("Evaluate", mock.AnythingOfType("*mocks.Context")).
								Return(types.False, nil)

							return expression
						}(),
						Command: NewSignedCommand("three-command"),
					},
				},
			},
			args: args{
				context: func() context.Context {
					context := new(contextmocks.Context)
					context.On("Copy").Return(context)

					return context
				}(),
			},
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name: "error/with condition evaluation",
			fields: fields{
				conditionalCases: []ConditionalCase{
					{
						Condition: func() Expression {
							expression := NewSignedExpression("one-condition")
							expression.
								On("Evaluate", mock.AnythingOfType("*mocks.Context")).
								Return(types.False, nil)

							return expression
						}(),
						Command: NewSignedCommand("one-command"),
					},
					{
						Condition: func() Expression {
							expression := NewSignedExpression("two-condition")
							expression.
								On("Evaluate", mock.AnythingOfType("*mocks.Context")).
								Return(nil, iotest.ErrTimeout)

							return expression
						}(),
						Command: NewSignedCommand("two-command"),
					},
					{
						Condition: NewSignedExpression("three-condition"),
						Command:   NewSignedCommand("three-command"),
					},
				},
			},
			args: args{
				context: func() context.Context {
					context := new(contextmocks.Context)
					context.On("Copy").Return(context)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error/with condition conversion",
			fields: fields{
				conditionalCases: []ConditionalCase{
					{
						Condition: func() Expression {
							expression := NewSignedExpression("one-condition")
							expression.
								On("Evaluate", mock.AnythingOfType("*mocks.Context")).
								Return(types.False, nil)

							return expression
						}(),
						Command: NewSignedCommand("one-command"),
					},
					{
						Condition: func() Expression {
							expression := NewSignedExpression("two-condition")
							expression.
								On("Evaluate", mock.AnythingOfType("*mocks.Context")).
								Return(func() {}, nil)

							return expression
						}(),
						Command: NewSignedCommand("two-command"),
					},
					{
						Condition: NewSignedExpression("three-condition"),
						Command:   NewSignedCommand("three-command"),
					},
				},
			},
			args: args{
				context: func() context.Context {
					context := new(contextmocks.Context)
					context.On("Copy").Return(context)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error/with command running",
			fields: fields{
				conditionalCases: []ConditionalCase{
					{
						Condition: func() Expression {
							expression := NewSignedExpression("one-condition")
							expression.
								On("Evaluate", mock.AnythingOfType("*mocks.Context")).
								Return(types.False, nil)

							return expression
						}(),
						Command: NewSignedCommand("one-command"),
					},
					{
						Condition: func() Expression {
							expression := NewSignedExpression("two-condition")
							expression.
								On("Evaluate", mock.AnythingOfType("*mocks.Context")).
								Return(types.True, nil)

							return expression
						}(),
						Command: func() runtime.Command {
							command := NewSignedCommand("two-command")
							command.On("Run", mock.AnythingOfType("*mocks.Context")).Return(nil, iotest.ErrTimeout)

							return command
						}(),
					},
					{
						Condition: NewSignedExpression("three-condition"),
						Command:   NewSignedCommand("three-command"),
					},
				},
			},
			args: args{
				context: func() context.Context {
					context := new(contextmocks.Context)
					context.On("Copy").Return(context)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			expression := ConditionalExpression{
				conditionalCases: data.fields.conditionalCases,
			}
			gotResult, gotErr := expression.Evaluate(data.args.context)

			checkConditionalCases(test, data.fields.conditionalCases)
			mock.AssertExpectationsForObjects(test, data.args.context)
			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}

func checkConditionalCases(test *testing.T, conditionalCases []ConditionalCase) {
	for _, conditionalCase := range conditionalCases {
		mock.AssertExpectationsForObjects(test, conditionalCase.Condition, conditionalCase.Command)
	}
}
