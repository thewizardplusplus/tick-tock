package commands

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	expressionsmocks "github.com/thewizardplusplus/tick-tock/runtime/expressions/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestSetCommand(test *testing.T) {
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
					state := context.State{Name: "test"}

					context := new(contextmocks.Context)
					context.On("SetState", state).Return(nil)

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
					expressionOne := new(expressionsmocks.Expression)
					expressionOne.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

					expressionTwo := new(expressionsmocks.Expression)
					expressionTwo.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(4.2, nil)

					return []expressions.Expression{expressionOne, expressionTwo}
				}(),
			},
			args: args{
				context: func() context.Context {
					state := context.State{
						Name:      "test",
						Arguments: []interface{}{2.3, 4.2},
					}

					context := new(contextmocks.Context)
					context.On("SetState", state).Return(nil)

					return context
				}(),
			},
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name: "error with argument evaluation",
			fields: fields{
				name: "test",
				arguments: func() []expressions.Expression {
					expressionOne := new(expressionsmocks.Expression)
					expressionOne.
						On("Evaluate", mock.AnythingOfType("*mocks.Context")).
						Return(nil, iotest.ErrTimeout)

					expressionTwo := new(expressionsmocks.Expression)

					return []expressions.Expression{expressionOne, expressionTwo}
				}(),
			},
			args: args{
				context: new(contextmocks.Context),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with state setting",
			fields: fields{
				name: "test",
				arguments: func() []expressions.Expression {
					expressionOne := new(expressionsmocks.Expression)
					expressionOne.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

					expressionTwo := new(expressionsmocks.Expression)
					expressionTwo.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(4.2, nil)

					return []expressions.Expression{expressionOne, expressionTwo}
				}(),
			},
			args: args{
				context: func() context.Context {
					state := context.State{
						Name:      "test",
						Arguments: []interface{}{2.3, 4.2},
					}

					context := new(contextmocks.Context)
					context.On("SetState", state).Return(iotest.ErrTimeout)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			gotResult, gotErr := NewSetCommand(testData.fields.name, testData.fields.arguments).
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
