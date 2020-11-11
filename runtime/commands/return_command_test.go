package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime"
)

func TestReturnCommand(test *testing.T) {
	context := new(MockContext)
	gotResult, gotErr := ReturnCommand{}.Run(context)

	mock.AssertExpectationsForObjects(test, context)
	assert.Nil(test, gotResult)
	assert.Equal(test, runtime.ErrReturn, gotErr)
}
