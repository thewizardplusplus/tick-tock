package commands

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	expressionsmocks "github.com/thewizardplusplus/tick-tock/runtime/expressions/mocks"
)

func TestExpressionCommand(test *testing.T) {
	type fields struct {
		expression expressions.Expression
	}

	for _, testData := range []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				expression: func() expressions.Expression {
					expression := new(expressionsmocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

					return expression
				}(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				expression: func() expressions.Expression {
					expression := new(expressionsmocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(nil, iotest.ErrTimeout)

					return expression
				}(),
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			context := new(contextmocks.Context)
			err := NewExpressionCommand(testData.fields.expression).Run(context)

			mock.AssertExpectationsForObjects(test, testData.fields.expression, context)
			testData.wantErr(test, err)
		})
	}
}
