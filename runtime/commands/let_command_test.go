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
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				identifier: "test",
				expression: func() expressions.Expression {
					expression := new(expressionsmocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

					return expression
				}(),
			},
			args: args{
				context: func() context.Context {
					context := new(contextmocks.Context)
					context.On("SetValue", "test", 2.3).Return()

					return context
				}(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				identifier: "test",
				expression: func() expressions.Expression {
					expression := new(expressionsmocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(nil, iotest.ErrTimeout)

					return expression
				}(),
			},
			args: args{
				context: new(contextmocks.Context),
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			err := NewLetCommand(testData.fields.identifier, testData.fields.expression).
				Run(testData.args.context)

			mock.AssertExpectationsForObjects(test, testData.fields.expression, testData.args.context)
			testData.wantErr(test, err)
		})
	}
}
