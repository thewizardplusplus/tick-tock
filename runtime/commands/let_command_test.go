package commands

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

func TestLetCommand(test *testing.T) {
	type fields struct {
		identifier string
		expression expressions.Expression
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
			name: "success",
			fields: fields{
				identifier: "test",
				expression: func() expressions.Expression {
					expression := new(MockExpression)
					expression.On("Evaluate", mock.AnythingOfType("*commands.MockContext")).Return(2.3, nil)

					return expression
				}(),
			},
			args: args{
				context: func() context.Context {
					context := new(MockContext)
					context.On("SetValue", "test", 2.3).Return()

					return context
				}(),
			},
			wantResult: 2.3,
			wantErr:    assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				identifier: "test",
				expression: func() expressions.Expression {
					expression := new(MockExpression)
					expression.On("Evaluate", mock.AnythingOfType("*commands.MockContext")).Return(nil, iotest.ErrTimeout)

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
		test.Run(testData.name, func(test *testing.T) {
			gotResult, gotErr := NewLetCommand(testData.fields.identifier, testData.fields.expression).
				Run(testData.args.context)

			mock.AssertExpectationsForObjects(test, testData.fields.expression, testData.args.context)
			assert.Equal(test, testData.wantResult, gotResult)
			testData.wantErr(test, gotErr)
		})
	}
}
