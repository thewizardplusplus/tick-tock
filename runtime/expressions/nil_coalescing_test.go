package expressions

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestNewNilCoalescingOperator(test *testing.T) {
	leftOperand := NewSignedExpression("one")
	rightOperand := NewSignedExpression("two")
	got := NewNilCoalescingOperator(leftOperand, rightOperand)

	mock.AssertExpectationsForObjects(test, leftOperand, rightOperand)
	assert.Equal(test, leftOperand, got.leftOperand)
	assert.Equal(test, rightOperand, got.rightOperand)
}

func TestNilCoalescingOperator_Evaluate(test *testing.T) {
	type fields struct {
		leftOperand  Expression
		rightOperand Expression
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
			name: "success/with the not nil left operand",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("one")
					expression.
						On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).
						Return(&types.Pair{Head: "one", Tail: &types.Pair{Head: "two", Tail: nil}}, nil)

					return expression
				}(),
				rightOperand: NewSignedExpression("two"),
			},
			args: args{
				context: new(MockContext),
			},
			wantResult: &types.Pair{Head: "one", Tail: &types.Pair{Head: "two", Tail: nil}},
			wantErr:    assert.NoError,
		},
		{
			name: "success/with the nil left operand",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("one")
					expression.
						On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).
						Return(types.Nil{}, nil)

					return expression
				}(),
				rightOperand: func() Expression {
					expression := NewSignedExpression("two")
					expression.
						On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).
						Return(&types.Pair{Head: "three", Tail: &types.Pair{Head: "four", Tail: nil}}, nil)

					return expression
				}(),
			},
			args: args{
				context: new(MockContext),
			},
			wantResult: &types.Pair{Head: "three", Tail: &types.Pair{Head: "four", Tail: nil}},
			wantErr:    assert.NoError,
		},
		{
			name: "error/with evaluation of the left operand",
			fields: fields{
				leftOperand: func() Expression {
					expression := NewSignedExpression("one")
					expression.
						On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).
						Return(nil, iotest.ErrTimeout)

					return expression
				}(),
				rightOperand: NewSignedExpression("two"),
			},
			args: args{
				context: new(MockContext),
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
						On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).
						Return(types.Nil{}, nil)

					return expression
				}(),
				rightOperand: func() Expression {
					expression := NewSignedExpression("two")
					expression.
						On("Evaluate", mock.AnythingOfType("*expressions.MockContext")).
						Return(nil, iotest.ErrTimeout)

					return expression
				}(),
			},
			args: args{
				context: new(MockContext),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			expression := NilCoalescingOperator{
				leftOperand:  data.fields.leftOperand,
				rightOperand: data.fields.rightOperand,
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
