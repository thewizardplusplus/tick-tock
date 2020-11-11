package commands

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestSendCommand(test *testing.T) {
	type fields struct {
		name      string
		arguments []expressions.Expression
	}
	type args struct {
		context context.Context
	}

	for _, testData := range []struct {
		name       string
		fields     fields
		args       args
		wantResult interface{}
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success without arguments",
			fields: fields{
				name:      "test",
				arguments: nil,
			},
			args: args{
				context: func() context.Context {
					message := context.Message{Name: "test"}

					context := new(MockContext)
					context.On("SendMessage", message).Return()

					return context
				}(),
			},
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name: "success with arguments",
			fields: fields{
				name: "test",
				arguments: func() []expressions.Expression {
					expressionOne := new(MockExpression)
					expressionOne.On("Evaluate", mock.AnythingOfType("*commands.MockContext")).Return(2.3, nil)

					expressionTwo := new(MockExpression)
					expressionTwo.On("Evaluate", mock.AnythingOfType("*commands.MockContext")).Return(4.2, nil)

					return []expressions.Expression{expressionOne, expressionTwo}
				}(),
			},
			args: args{
				context: func() context.Context {
					message := context.Message{
						Name:      "test",
						Arguments: []interface{}{2.3, 4.2},
					}

					context := new(MockContext)
					context.On("SendMessage", message).Return()

					return context
				}(),
			},
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				name: "test",
				arguments: func() []expressions.Expression {
					expressionOne := new(MockExpression)
					expressionOne.
						On("Evaluate", mock.AnythingOfType("*commands.MockContext")).
						Return(nil, iotest.ErrTimeout)

					expressionTwo := new(MockExpression)

					return []expressions.Expression{expressionOne, expressionTwo}
				}(),
			},
			args: args{
				context: new(MockContext),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			gotResult, gotErr := NewSendCommand(testData.fields.name, testData.fields.arguments).
				Run(testData.args.context)

			mock.AssertExpectationsForObjects(test, testData.args.context)
			for _, argument := range testData.fields.arguments {
				mock.AssertExpectationsForObjects(test, argument)
			}
			assert.Equal(test, testData.wantResult, gotResult)
			testData.wantErr(test, gotErr)
		})
	}
}
