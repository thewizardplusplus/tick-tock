package commands

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	expressionsmocks "github.com/thewizardplusplus/tick-tock/runtime/expressions/mocks"
)

func TestExpressionCommand(test *testing.T) {
	type fields struct {
		expression expressions.Expression
	}

	for _, testData := range []struct {
		name       string
		fields     fields
		wantResult interface{}
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				expression: func() expressions.Expression {
					expression := new(expressionsmocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*commands.MockContext")).Return(2.3, nil)

					return expression
				}(),
			},
			wantResult: 2.3,
			wantErr:    assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				expression: func() expressions.Expression {
					expression := new(expressionsmocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*commands.MockContext")).Return(nil, iotest.ErrTimeout)

					return expression
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			context := new(MockContext)
			gotResult, gotErr := NewExpressionCommand(testData.fields.expression).Run(context)

			mock.AssertExpectationsForObjects(test, testData.fields.expression, context)
			assert.Equal(test, testData.wantResult, gotResult)
			testData.wantErr(test, gotErr)
		})
	}
}
