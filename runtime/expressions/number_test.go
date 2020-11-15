package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewNumber(test *testing.T) {
	got := NewNumber(2.3)

	assert.Equal(test, 2.3, got.value)
}

func TestNumber_Evaluate(test *testing.T) {
	context := new(MockContext)
	number := Number{2.3}
	gotResult, gotErr := number.Evaluate(context)

	mock.AssertExpectationsForObjects(test, context)
	assert.Equal(test, 2.3, gotResult)
	assert.NoError(test, gotErr)
}
