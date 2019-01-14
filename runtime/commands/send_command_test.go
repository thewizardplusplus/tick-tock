package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestSendCommand(test *testing.T) {
	context := new(mocks.Context)
	context.On("SendMessage", "test").Return()

	err := NewSendCommand("test").Run(context)

	context.AssertExpectations(test)
	assert.NoError(test, err)
}
