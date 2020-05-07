package expressions

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestNewBooleanOperator(test *testing.T) {
	leftOperand := NewSignedExpression("one")
	rightOperand := NewSignedExpression("two")
	got := NewBooleanOperator(leftOperand, rightOperand, types.True)

	mock.AssertExpectationsForObjects(test, leftOperand, rightOperand)
	assert.Equal(test, leftOperand, got.leftOperand)
	assert.Equal(test, rightOperand, got.rightOperand)
	assert.Equal(test, types.True, got.valueForEarlyExit)
}

func TestBooleanOperator_Evaluate(test *testing.T) {
	type fields struct {
		leftOperand       Expression
		rightOperand      Expression
		valueForEarlyExit types.Boolean
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
			name: "success/with early exit on true/with the true left operand",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("one")
					expression.
						On("Evaluate", mock.AnythingOfType("*mocks.Context")).
						Return(&types.Pair{Head: "one", Tail: &types.Pair{Head: "two", Tail: nil}}, nil)

					return expression
				}(),
				rightOperand:      NewSignedExpression("two"),
				valueForEarlyExit: types.True,
			},
			args: args{
				context: new(contextmocks.Context),
			},
			wantResult: &types.Pair{Head: "one", Tail: &types.Pair{Head: "two", Tail: nil}},
			wantErr:    assert.NoError,
		},
		{
			name: "success/with early exit on true/with the false left operand",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("one")
					expression.
						On("Evaluate", mock.AnythingOfType("*mocks.Context")).
						Return((*types.Pair)(nil), nil)

					return expression
				}(),
				rightOperand: func() Expression {
					expression := NewSignedExpression("two")
					expression.
						On("Evaluate", mock.AnythingOfType("*mocks.Context")).
						Return(&types.Pair{Head: "three", Tail: &types.Pair{Head: "four", Tail: nil}}, nil)

					return expression
				}(),
				valueForEarlyExit: types.True,
			},
			args: args{
				context: new(contextmocks.Context),
			},
			wantResult: &types.Pair{Head: "three", Tail: &types.Pair{Head: "four", Tail: nil}},
			wantErr:    assert.NoError,
		},
		{
			name: "success/with early exit on false/with the true left operand",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("one")
					expression.
						On("Evaluate", mock.AnythingOfType("*mocks.Context")).
						Return(&types.Pair{Head: "one", Tail: &types.Pair{Head: "two", Tail: nil}}, nil)

					return expression
				}(),
				rightOperand: func() Expression {
					expression := NewSignedExpression("two")
					expression.
						On("Evaluate", mock.AnythingOfType("*mocks.Context")).
						Return(&types.Pair{Head: "three", Tail: &types.Pair{Head: "four", Tail: nil}}, nil)

					return expression
				}(),
				valueForEarlyExit: types.False,
			},
			args: args{
				context: new(contextmocks.Context),
			},
			wantResult: &types.Pair{Head: "three", Tail: &types.Pair{Head: "four", Tail: nil}},
			wantErr:    assert.NoError,
		},
		{
			name: "success/with early exit on false/with the false left operand",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("one")
					expression.
						On("Evaluate", mock.AnythingOfType("*mocks.Context")).
						Return((*types.Pair)(nil), nil)

					return expression
				}(),
				rightOperand:      NewSignedExpression("two"),
				valueForEarlyExit: types.False,
			},
			args: args{
				context: new(contextmocks.Context),
			},
			wantResult: (*types.Pair)(nil),
			wantErr:    assert.NoError,
		},
		{
			name: "error/with evaluation of the left operand",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("one")
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(nil, iotest.ErrTimeout)

					return expression
				}(),
				rightOperand:      NewSignedExpression("two"),
				valueForEarlyExit: types.True,
			},
			args: args{
				context: new(contextmocks.Context),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error/with conversion of the left operand",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("one")
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(func() {}, nil)

					return expression
				}(),
				rightOperand:      NewSignedExpression("two"),
				valueForEarlyExit: types.True,
			},
			args: args{
				context: new(contextmocks.Context),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error/with evaluation of the right operand",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("one")
					expression.
						On("Evaluate", mock.AnythingOfType("*mocks.Context")).
						Return((*types.Pair)(nil), nil)

					return expression
				}(),
				rightOperand: func() Expression {
					expression := NewSignedExpression("two")
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(nil, iotest.ErrTimeout)

					return expression
				}(),
				valueForEarlyExit: types.True,
			},
			args: args{
				context: new(contextmocks.Context),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			expression := BooleanOperator{
				leftOperand:       data.fields.leftOperand,
				rightOperand:      data.fields.rightOperand,
				valueForEarlyExit: data.fields.valueForEarlyExit,
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
