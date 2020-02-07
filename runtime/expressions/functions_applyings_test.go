package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestNewMultiplication(test *testing.T) {
	leftOperand := NewSignedExpression("left")
	leftOperand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.0, nil)

	rightOperand := NewSignedExpression("right")
	rightOperand.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(3.0, nil)

	context := new(contextmocks.Context)
	expression := NewMultiplication(leftOperand, rightOperand)
	gotResult, gotErr := expression.Evaluate(context)

	mock.AssertExpectationsForObjects(test, leftOperand, rightOperand, context)
	assert.Equal(test, 6.0, gotResult)
	assert.NoError(test, gotErr)
}
